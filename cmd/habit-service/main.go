package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/glados"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/syslog"
)

var (
	cfg    config.Config
	client habit.Client
	log    syslog.Logger
)

func init() {
	var err error
	cfg, err = config.ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Could not read config variables: %v", err)
		os.Exit(1)
	}

	log = syslog.NewLogger(cfg.Telegram.ChatId, cfg.Telegram.Token)

	location, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		log.Warn(
			"Warning: invalid timezone location: '%s', falling back to UTC: %v",
			cfg.TimezoneLocation,
			err,
		)
		location, _ = time.LoadLocation("UTC")
	}

	client, err = habit.GetClient(context.Background(), cfg.Habit.SpreadsheetId, location)
	if err != nil {
		log.Fatal("Could not create gsheets client for Habit Service: %v", err)
	}
}

func main() {
	// start HTTP server
	http.HandleFunc("/", getDueTasks)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil)

	// start Glados command handler
	listener := glados.NewListener(cfg.Glados, log)
	go listener.Listen(context.Background(), handleGladosCommand)

	// shutdown gracefully upon termination
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	listener.Close()
	log.Info("Shutting down Habit Service")
}

func getDueTasks(w http.ResponseWriter, req *http.Request) {
	cards, err := client.FetchNewCards()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		message := fmt.Sprintf("could not fetch new cards: %v", err)
		log.Info(message)
		fmt.Fprintf(w, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

// @todo: create a command parser if there are more commands to handle in the future
func handleGladosCommand(args []string) {
	if len(args) != 2 || args[0] != "mark" {
		log.Error("Could not parse Glados command from args: %v", args)
		return
	}

	// @todo: mark habit
}
