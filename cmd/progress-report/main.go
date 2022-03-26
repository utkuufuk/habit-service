package main

import (
	"context"
	"os"

	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
)

func main() {
	ctx := context.Background()
	client, err := habit.GetClient(ctx)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	_, err = service.ReportProgressAction{}.Run(ctx, client)
	if err != nil {
		logger.Error("Could not run Glados command: %v", err)
		os.Exit(1)
	}
}
