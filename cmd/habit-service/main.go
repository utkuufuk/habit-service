package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/glados"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/service"
)

var (
	cfg    config.Config
	client habit.Client
)

func main() {
	var err error
	cfg, err = config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Could not read config variables: %v", err)
	}

	client, err = habit.GetClient(context.Background(), cfg.SpreadsheetId, cfg.TimezoneLocation)
	if err != nil {
		log.Fatalf("Could not create gsheets client for Habit Service: %v", err)
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

	action := service.FetchHabitsAsTrelloCardsAction{
		TimezoneLocation: cfg.TimezoneLocation,
	}
	cards, err := action.Run(req.Context(), client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("could not fetch new cards: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func handleGladosCommand(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type response struct {
		Message string `json:"message"`
	}
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{fmt.Sprintf("Could not read HTTP request body: %v", err)})
		return
	}

	var request struct {
		Args []string `json:"args"`
	}
	if err = json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{fmt.Sprintf("Could not decode HTTP request body: %v", err)})
		return
	}

	action, err := glados.ParseCommand(request.Args, cfg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{fmt.Sprintf("Could not parse Glados command: %v", err)})
		return
	}

	message, err := action.Run(req.Context(), client)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{err.Error()})
		return
	}

	json.NewEncoder(w).Encode(response{message})
}
