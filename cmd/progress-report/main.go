package main

import (
	"context"
	"os"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
)

func main() {
	cfg, err := config.ParseProgressReportConfig()
	if err != nil {
		logger.Error("Failed to parse server config: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := habit.GetClient(ctx, cfg.GoogleSheets)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	action := service.ReportProgressAction{
		TelegramChatId:   cfg.TelegramChatId,
		TelegramToken:    cfg.TelegramToken,
		TimezoneLocation: cfg.TimezoneLocation,
	}

	_, err = action.Run(ctx, client)
	if err != nil {
		logger.Error("Could not run Glados command: %v", err)
		os.Exit(1)
	}
}
