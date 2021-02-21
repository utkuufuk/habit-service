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
	cards, err := client.FetchNewCards(req.Context(), time.Now().In(loc))
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
