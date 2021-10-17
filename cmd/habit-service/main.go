package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/entrello"
	"github.com/utkuufuk/habit-service/internal/glados"
	"github.com/utkuufuk/habit-service/internal/habit"
)

var (
	cfg      config.Config
	client   habit.Client
	location *time.Location
)

func main() {
	var err error
	cfg, err = config.ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Could not read config variables: %v", err)
		os.Exit(1)
	}

	location, err = time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		log.Fatalf("Invalid timezone location: '%s': %v", cfg.TimezoneLocation, err)
	}

	client, err = habit.GetClient(context.Background(), cfg.SpreadsheetId, location)
	if err != nil {
		log.Fatalf("Could not create gsheets client for Habit Service: %v", err)
	}

	msg, err := glados.RunCommand(context.Background(), client, location, []string{}, cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(msg)
	return

	// http.HandleFunc("/entrello", handleEntrelloRequest)
	// http.HandleFunc("/glados", handleGladosCommand)
	// http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), nil)
}

func handleEntrelloRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cards, err := entrello.FetchHabitCards(client, location)
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

	message, err := glados.RunCommand(context.Background(), client, location, request.Args, cfg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{err.Error()})
		return
	}

	json.NewEncoder(w).Encode(response{message})
}
