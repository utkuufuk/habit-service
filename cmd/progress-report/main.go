package main

import (
	"context"
	"log"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/service"
)

func main() {
	ctx := context.Background()
	client, err := habit.GetClient(ctx, config.Values.SpreadsheetId)
	if err != nil {
		log.Fatalf("Could not create gsheets client for Habit Service: %v", err)
	}

	action := service.ReportProgressAction{
		TelegramChatId:   config.Values.TelegramChatId,
		TelegramToken:    config.Values.TelegramToken,
		TimezoneLocation: config.Values.TimezoneLocation,
	}
	_, err = action.Run(ctx, client)
	if err != nil {
		log.Fatalf("Could not run Glados command: %v", err)
	}
}
