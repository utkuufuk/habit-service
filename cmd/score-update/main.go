package main

import (
	"context"
	"log"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
)

func main() {
	client, err := habit.GetClient(context.Background(), config.Values.SpreadsheetId)
	if err != nil {
		log.Fatalf("Could not create gsheets client for Habit Service: %v", err)
	}

	now := time.Now().In(config.Values.TimezoneLocation)
	habits, err := client.FetchHabits(now)
	if err != nil {
		log.Fatalf("could not fetch habits: %v", err)
	}

	if err = client.UpdateScores(habits, now); err != nil {
		log.Fatalf("could not update habit scores: %v", err)
	}
}
