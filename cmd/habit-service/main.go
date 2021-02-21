package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
)

var (
	cfg    config.Config
	client habit.Client
	loc    *time.Location
)

func getDueTasks(w http.ResponseWriter, req *http.Request) {
	var params struct {
		Label string `json:"label"`
	}

	if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
		fmt.Printf("could not decode label: %v", err)
	}

	cards, err := client.FetchNewCards(req.Context(), time.Now().In(loc), params.Label)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not fetch new cards: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func main() {
	var err error
	cfg, err = config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Could not read config variables: %v", err)
	}

	loc, err = time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		log.Fatalf("invalid timezone location: %v", loc)
	}

	client = habit.GetClient(cfg.Habit.SpreadsheetId)

	http.HandleFunc("/", getDueTasks)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil)
}
