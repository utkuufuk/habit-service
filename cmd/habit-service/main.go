package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
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
	http.HandleFunc("/", getDueTasks)
	http.HandleFunc("/mark", markHabit)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil)
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

func markHabit(w http.ResponseWriter, req *http.Request) {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Could not parse HTTP request body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	var body struct {
		Cell   string `json:"cell"`
		Symbol string `json:"symbol"`
	}
	json.Unmarshal(reqBody, &body)
	cell := body.Cell
	symbol := body.Symbol

	matched, err := regexp.MatchString(`[a-zA-Z]{3}\ 202\d\![A-Z][1-9][0-9]?$|^100$`, cell)
	if err != nil || matched == false {
		log.Error("Invalid cell '%s' to mark habit in Glados command: %v", cell, err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if !habit.IsValidMarkSymbol(symbol) {
		log.Error("Invalid symbol '%s' to mark habit in HTTP request", symbol)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := client.MarkHabit(cell, symbol); err != nil {
		log.Error("Could not mark Habit on cell '%s' with symbol '%s': %v", cell, symbol, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
