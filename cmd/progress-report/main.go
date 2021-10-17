package main

import (
	"context"
	"log"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/service"
)

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Could not read config variables: %v", err)
	}

	ctx := context.Background()
	client, err := habit.GetClient(ctx, cfg.SpreadsheetId, cfg.TimezoneLocation)
	if err != nil {
		log.Fatalf("Could not create gsheets client for Habit Service: %v", err)
	}

	action := service.ReportProgressAction{
		TelegramChatId:   cfg.TelegramChatId,
		TelegramToken:    cfg.TelegramToken,
		TimezoneLocation: cfg.TimezoneLocation,
	}
	message, err := action.Run(ctx, client)
	if err != nil {
		log.Fatalf("Could not run Glados command: %v", err)
	}

	log.Println(message)
}
