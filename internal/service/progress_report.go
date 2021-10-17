package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/tableimage"
)

type ReportProgressAction struct {
	TelegramChatId   int64
	TelegramToken    string
	TimezoneLocation *time.Location
}

func (a ReportProgressAction) Run(ctx context.Context, client habit.Client) (string, error) {
	now := time.Now().In(a.TimezoneLocation)
	currentHabits, err := client.FetchHabits(now)
	if err != nil {
		return "", fmt.Errorf("could not fetch this month's habits: %v\n", err)
	}

	year, month, _ := now.Date()
	lastMonth := time.Date(year, month, 1, 0, 0, 0, 0, a.TimezoneLocation).Add(-time.Nanosecond)
	previousHabits, err := client.FetchHabits(lastMonth)
	if err != nil {
		return "", fmt.Errorf("could not fetch habits from last month: %v\n", err)
	}

	table := tableimage.NewTable()
	for name, habit := range currentHabits {
		table.AddRow(name, previousHabits[name].Score*100, habit.Score*100)
	}

	path := fmt.Sprintf("./progress-reports/%s.png", now.Format("2006-01-02T15:04:05"))
	table.Save(path)
	return "habit progress report sent", a.sendProgressReport(path)
}

func (a ReportProgressAction) sendProgressReport(path string) error {
	bot, err := tgbotapi.NewBotAPI(a.TelegramToken)
	if err != nil {
		return fmt.Errorf("could not initialize Telegram bot client: %w", err)
	}

	photoBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read progress report image '%s': %w", path, err)
	}

	_, err = bot.Send(tgbotapi.NewPhotoUpload(a.TelegramChatId, tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: photoBytes,
	}))
	return err
}
