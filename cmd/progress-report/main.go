package main

import (
	"context"
	"os"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
	"github.com/utkuufuk/habit-service/internal/sheets"
)

func main() {
	cfg, err := config.ParseProgressReportConfig()
	if err != nil {
		logger.Error("Failed to parse server config: %v", err)
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := sheets.GetClient(ctx, cfg.GoogleSheets)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	if err = service.ReportProgress(
		client,
		cfg.TimezoneLocation,
		cfg.SkipList,
		cfg.TelegramChatId,
		cfg.TelegramToken,
	); err != nil {
		logger.Error("Could not run Glados command: %v", err)
		os.Exit(1)
	}
}
