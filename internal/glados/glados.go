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

	fmt.Println("This Month:")
	for _, habit := range currentHabits {
		fmt.Printf("%s: %s%%\n", habit.Name, strconv.FormatFloat(habit.Score*100, 'f', 0, 32))
	}

	y, m, _ := now.Date()
	for _, i := range []int{1, 2, 3} {
		month := time.Month(int(m) + 1 - i)
		lastMonth := time.Date(y, month, 1, 0, 0, 0, 0, location).Add(-time.Nanosecond)
		habits, err := client.FetchHabits(lastMonth)
		if err != nil {
			return "", fmt.Errorf("could not fetch habits from %v: %v\n", (month - 1), err)
		}

		fmt.Printf("\n%v:\n", month-1)
		for _, habit := range habits {
			fmt.Printf("%s: %s%%\n", habit.Name, strconv.FormatFloat(habit.Score*100, 'f', 0, 32))
		}
	}

	// @todo
	return "", nil
}

func markHabit(client habit.Client, cell, symbol string) error {
	matched, err := regexp.MatchString(`[a-zA-Z]{3}\ 202\d\![A-Z][1-9][0-9]?$|^100$`, cell)
	if err != nil || matched == false {
		return fmt.Errorf("Invalid cell '%s' to mark habit in Glados command: %v", cell, err)
	}

	if !habit.IsValidMarkSymbol(symbol) {
		return fmt.Errorf("Invalid symbol '%s' to mark habit in HTTP request", symbol)
	}

	if err := client.MarkHabit(cell, symbol); err != nil {
		return fmt.Errorf("Could not mark Habit on cell '%s' with symbol '%s': %v", cell, symbol, err)
	}

	return nil
}
