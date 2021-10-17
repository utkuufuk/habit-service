package glados

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
)

func reportProgress(
	ctx context.Context,
	client habit.Client,
	location *time.Location,
	mailgunCfg config.Mailgun,
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
	return "habit progress report email sent", emailProgressReport(ctx, mailgunCfg, path, now)
}

func emailProgressReport(ctx context.Context, mailgunCfg config.Mailgun, path string, now time.Time) error {
	mg := mailgun.NewMailgun(mailgunCfg.Domain, mailgunCfg.ApiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)
	m := mg.NewMessage(
		mailgunCfg.From,
		fmt.Sprintf("Habit Progress Report (%s)", now.Format("Jan-02-2006")),
		"",
		mailgunCfg.To,
	)
	m.AddAttachment(path)

	_, _, err := mg.Send(ctx, m)
	return err
}
