package main

import (
	"context"
	"os"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/logger"
	"github.com/utkuufuk/habit-service/internal/service"
	"github.com/utkuufuk/habit-service/internal/sheets"
)

func main() {
	loc, gSheetsCfg := config.ParseCommonConfig()

	ctx := context.Background()
	client, err := sheets.GetClient(ctx, gSheetsCfg)
	if err != nil {
		logger.Error("Could not create gsheets client for Habit Service: %v", err)
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		displayCards(client, loc)
		return
	}

	switch os.Args[1] {
	case "progress-report":
		sendPorgressReport(client)
	default:
		logger.Error("Uncrecognized command: %s", os.Args[1])
	}
}

func displayCards(client sheets.Client, loc *time.Location) {
	cards, err := service.FetchHabitCards(client, loc)
	if err != nil {
		logger.Error("Could not fetch undone habits: %v", err)
		return
	}
	for _, card := range cards {
		logger.Info("%s / %s / %s", card.Name, card.Desc, card.Due.Format("2006-01-02"))
	}
}

func sendPorgressReport(client sheets.Client) {
	cfg, err := config.ParseProgressReportConfig()
	if err != nil {
		logger.Error("Failed to parse progress report config: %v", err)
		os.Exit(1)
	}

	if err = service.ReportProgress(
		client,
		cfg.TimezoneLocation,
		cfg.SkipList,
		cfg.TelegramChatId,
		cfg.TelegramToken,
	); err != nil {
		logger.Error("Could not run send progress report: %v", err)
	}
}
