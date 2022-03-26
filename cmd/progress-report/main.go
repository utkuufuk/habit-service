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
	ctx := context.Background()
	client, err := habit.GetClient(ctx, config.SpreadsheetId)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	action := service.ReportProgressAction{
		TelegramChatId:   config.TelegramChatId,
		TelegramToken:    config.TelegramToken,
		TimezoneLocation: config.TimezoneLocation,
	}
	_, err = action.Run(ctx, client)
	if err != nil {
		logger.Error("Could not run Glados command: %v", err)
		os.Exit(1)
	}
}
