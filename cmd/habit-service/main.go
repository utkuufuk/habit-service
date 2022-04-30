package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
	"github.com/utkuufuk/habit-service/internal/sheets"
)

var (
	cfg    config.ServerConfig
	client sheets.Client
)

func init() {
	var err error
	cfg, err = config.ParseServerConfig()
	if err != nil {
		logger.Error("Failed to parse server config: %v", err)
		os.Exit(1)
	}

	client, err = sheets.GetClient(context.Background(), cfg.GoogleSheets)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/entrello", handleEntrelloRequest)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil)
}

func handleEntrelloRequest(w http.ResponseWriter, req *http.Request) {
	if cfg.Secret != "" && req.Header.Get("X-Api-Key") != cfg.Secret {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if req.Method == http.MethodGet {
		cards, err := service.FetchHabitCards(client, cfg.TimezoneLocation)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, fmt.Sprintf("could not fetch new cards: %v", err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
		return
	}

	if req.Method == http.MethodPost {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error("Could not read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var card service.TrelloCard
		cell := strings.Split(card.Desc, "\n")[0]
		if err = json.Unmarshal(body, &card); err != nil {
			logger.Warn("Invalid request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = service.UpdateHabit(client, cfg.TimezoneLocation, cell, card.Labels)
		if err != nil {
			logger.Error("Could not update habit at cell '%s': %v", cell, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
