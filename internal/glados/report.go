package glados

import (
	"context"
	"fmt"
	"strconv"
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
		sign := ""
		if habit.Score > previousHabits[name].Score {
			sign = "+"
		}

		currentScore := strconv.FormatFloat(habit.Score*100, 'f', 0, 32)
		lastScore := strconv.FormatFloat(previousHabits[name].Score*100, 'f', 0, 32)
		delta := strconv.FormatFloat((habit.Score-previousHabits[name].Score)*100, 'f', 0, 32)
		table.addRow(name, lastScore, currentScore, sign+delta)
	}

	path := fmt.Sprintf("./reports/progress-report-%s.png", now.Format("2006-01-02T15:04:05"))
	table.save(path)
	return "habit progress report email sent", sendProgressReportEmail(
		ctx,
		mailgunCfg.ApiKey,
		mailgunCfg.Domain,
		mailgunCfg.From,
		mailgunCfg.To,
		path,
	)
}

func sendProgressReportEmail(ctx context.Context, apiKey, domain, from, to, path string) error {
	mg := mailgun.NewMailgun(domain, apiKey)
	mg.SetAPIBase(mailgun.APIBaseEU)
	m := mg.NewMessage(from, "Habit Progress Report", "Testing some Mailgun awesomeness!", to)
	m.AddAttachment(path)

	_, _, err := mg.Send(ctx, m)
	return err
}
