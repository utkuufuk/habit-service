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
	"github.com/utkuufuk/habit-service/internal/glados"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
)

var (
	client habit.Client
)

func main() {
	var err error
	client, err = habit.GetClient(context.Background())
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	http.HandleFunc("/entrello", handleEntrelloRequest)
	// http.HandleFunc("/glados", handleGladosCommand)
	http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), nil)
}

func handleEntrelloRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		action := service.FetchHabitsAsTrelloCardsAction{}
		cards, err := action.Run(req.Context(), client)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, fmt.Sprintf("could not fetch new cards: %v", err))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
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
		// var card interface{}
		if err = json.Unmarshal(body, &card); err != nil {
			logger.Warn("Invalid request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// logger.Info("Card: %v", card)
		logger.Info("Desc: %s\nLabels: %v", card.Desc, card.Labels)

		cell := strings.Split(card.Desc, "\n")[0]
		matched, err := regexp.MatchString(`[a-zA-Z]{3} 202\d![A-Z][1-9][0-9]?$|^100$`, cell)
		if err != nil || matched == false {
			logger.Error("Invalid cell name '%s' in card description: %v", cell, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		logger.Info("Cell: %s", cell)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return
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

	action, err := glados.ParseCommand(request.Args)
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
