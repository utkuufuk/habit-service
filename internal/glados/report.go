package glados

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
)

func reportProgress(
	ctx context.Context,
	client habit.Client,
	location *time.Location,
	telegramCfg config.Telegram,
) (string, error) {
	now := time.Now().In(location)
	currentHabits, err := client.FetchHabits(now)
	if err != nil {
		return "", fmt.Errorf("could not fetch this month's habits: %v\n", err)
	}

	year, month, _ := now.Date()
	lastMonth := time.Date(year, month, 1, 0, 0, 0, 0, location).Add(-time.Nanosecond)
	previousHabits, err := client.FetchHabits(lastMonth)
	if err != nil {
		return "", fmt.Errorf("could not fetch habits from last month: %v\n", err)
	}

	table := newTable()
	for name, habit := range currentHabits {
		table.addRow(name, previousHabits[name].Score*100, habit.Score*100)
	}

	path := fmt.Sprintf("./reports/progress-report-%s.png", now.Format("2006-01-02T15:04:05"))
	table.save(path)
	return "habit progress report sent", sendProgressReport(telegramCfg, path)
}

func sendProgressReport(telegramCfg config.Telegram, path string) error {
	bot, err := tgbotapi.NewBotAPI(telegramCfg.Token)
	if err != nil {
		return fmt.Errorf("could not initialize Telegram bot client: %w", err)
	}

	photoBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read progress report image '%s': %w", path, err)
	}

	_, err = bot.Send(tgbotapi.NewPhotoUpload(telegramCfg.ChatId, tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: photoBytes,
	}))
	return err
}
