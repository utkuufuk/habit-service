package glados

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/utkuufuk/habit-service/internal/habit"
)

func RunCommand(client habit.Client, location *time.Location, args []string) (string, error) {
	if len(args) == 0 {
		return summarizeProgress(client, location)
	}

	if args[0] == "mark" && len(args) == 3 {
		return "", markHabit(client, args[0], args[1])
	}

	return "", fmt.Errorf("could not parse glados command from args: '%v'", args)
}

func summarizeProgress(client habit.Client, location *time.Location) (string, error) {
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

	message := ""
	for name, habit := range currentHabits {
		sign := ""
		if habit.Score > previousHabits[name].Score {
			sign = "+"
		}
		message += fmt.Sprintf(
			"%s:\n===============\nThis month: %s\nLast month: %s\nDelta: %s%s\n\n",
			name,
			strconv.FormatFloat(habit.Score*100, 'f', 0, 32),
			strconv.FormatFloat(previousHabits[name].Score*100, 'f', 0, 32),
			sign,
			strconv.FormatFloat((habit.Score-previousHabits[name].Score)*100, 'f', 0, 32),
		)
	}

	return message, nil
}

func markHabit(client habit.Client, cell, symbol string) error {
	matched, err := regexp.MatchString(`[a-zA-Z]{3}\ 202\d\![A-Z][1-9][0-9]?$|^100$`, cell)
	if err != nil || matched == false {
		return fmt.Errorf("invalid cell '%s' to mark habit in Glados command: %v", cell, err)
	}

	if !habit.IsValidMarkSymbol(symbol) {
		return fmt.Errorf("invalid habit symbol '%s' to mark habit", symbol)
	}

	if err := client.MarkHabit(cell, symbol); err != nil {
		return fmt.Errorf("could not mark habit on cell '%s' with symbol '%s': %v", cell, symbol, err)
	}

	return nil
}
