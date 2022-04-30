package service

import (
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/sheets"
)

const (
	dueHour = 23
)

func FetchEntrelloCards(client sheets.Client, loc *time.Location) ([]trello.Card, error) {
	now := time.Now().In(loc)

	habits, err := habit.FetchAll(client, now)
	if err != nil {
		return nil, fmt.Errorf("could not fetch habits: %w", err)
	}

	cards, err := toTrelloCards(habits, now)
	if err != nil {
		return nil, fmt.Errorf("could not convert habits to Trello cards: %w", err)
	}

	return cards, nil
}

// toTrelloCards returns a slice of trello cards from the given habits which haven't been marked today
func toTrelloCards(habits map[string]habit.Habit, now time.Time) (cards []trello.Card, err error) {
	for name, habit := range habits {
		if habit.State != "" {
			continue
		}

		// include the day of month in card title to force overwrite in the beginning of the next day
		title := fmt.Sprintf("%v (%d)", name, now.Day())

		due := time.Date(now.Year(), now.Month(), now.Day(), dueHour, 0, 0, 0, now.Location())
		c, err := trello.NewCard(title, string(habit.Cell), &due)
		if err != nil {
			return nil, fmt.Errorf("could not create habit card: %w", err)
		}

		cards = append(cards, c)
	}

	return cards, nil
}
