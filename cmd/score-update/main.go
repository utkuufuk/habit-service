package main

import (
	"context"
	"os"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/logger"
)

func main() {
	client, err := habit.GetClient(context.Background(), config.SpreadsheetId)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	now := time.Now().In(config.TimezoneLocation)
	habits, err := client.FetchHabits(now)
	if err != nil {
		logger.Error("could not fetch habits: %v", err)
		os.Exit(1)
	}

	if err = client.UpdateScores(habits, now); err != nil {
		logger.Error("could not update habit scores: %v", err)
		os.Exit(1)
	}
}
