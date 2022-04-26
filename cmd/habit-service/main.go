package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
)

var (
	cfg    config.ServerConfig
	client habit.Client
)

func init() {
	var err error
	cfg, err = config.ParseServerConfig()
	if err != nil {
		logger.Error("Failed to parse server config: %v", err)
		os.Exit(1)
	}

	client, err = habit.GetClient(context.Background(), cfg.GoogleSheets)
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
		return
	}

	if req.Method == http.MethodPost {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error("Could not read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var card struct {
			Desc   string `json:"desc"`
			Labels []struct {
				Name string `json:"name"`
			} `json:"labels"`
		}
		if err = json.Unmarshal(body, &card); err != nil {
			logger.Warn("Invalid request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cell := strings.Split(card.Desc, "\n")[0]
		matched, err := regexp.MatchString(`[a-zA-Z]{3} 202\d![A-Z][1-9][0-9]?$|^100$`, cell)
		if err != nil || matched == false {
			logger.Error("Invalid cell name '%s' in card description: %v", cell, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		symbol := "✔"
		for _, c := range card.Labels {
			if c.Name == "habit-skip" {
				symbol = "–"
				break
			}
			if c.Name == "habit-fail" {
				symbol = "✘"
				break
			}
		}

		_, err = service.MarkHabitAction{Cell: cell, Symbol: symbol}.Run(req.Context(), client)
		if err != nil {
			logger.Error("Could not mark habit at cell '%s' as %s: %v", cell, symbol, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
