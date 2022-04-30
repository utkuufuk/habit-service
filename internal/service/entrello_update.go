package service

import (
	"fmt"
	"time"

	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/sheets"
	"golang.org/x/exp/maps"
)

// UpdateHabit marks the habit as done/failed/skipped,
// and updates the scores of all habits in the spreadsheet.
func UpdateHabit(client sheets.Client, loc *time.Location, cell, symbol string) error {
	if err := habit.Mark(client, cell, symbol); err != nil {
		return fmt.Errorf("could not mark habit at cell '%s' with %s: %w", cell, symbol, err)
	}

	now := time.Now().In(loc)
	habits, err := habit.FetchAll(client, now)
	if err != nil {
		return fmt.Errorf("could not fetch habits: %w", err)
	}

	if err = habit.WriteScores(client, now, maps.Values(habits)); err != nil {
		return fmt.Errorf("could not write habit scores: %w", err)
	}

	return nil
}
