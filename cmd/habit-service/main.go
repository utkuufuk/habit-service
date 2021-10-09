package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/syslog"
)

var (
	client   habit.Client
	location *time.Location
)

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Could not read config variables: %v", err)
		os.Exit(1)
	}

	log := syslog.NewLogger(cfg.Telegram.ChatId, cfg.Telegram.Token)

	location, err = time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		log.Fatal("Invalid timezone location: '%s': %v", err)
	}

	client, err = habit.GetClient(context.Background(), cfg.Habit.SpreadsheetId, location)
	if err != nil {
		log.Fatal("Could not create gsheets client for Habit Service: %v", err)
	}

	http.HandleFunc("/entrello", handleEntrelloRequest)
	http.HandleFunc("/glados", handleGladosCommand)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil)
}

func handleEntrelloRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	habits, err := fetchHabitsForEntrello()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("could not fetch new cards: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

func handleGladosCommand(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(w, "Could not parse HTTP request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var response struct {
		Args []string `json:"args"`
	}
	json.Unmarshal(body, &response)

	message, err := runGladosCommand(response.Args)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, message)
	w.WriteHeader(http.StatusOK)
}
