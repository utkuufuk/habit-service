package entrello

import (
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/habit-service/internal/habit"
)

const (
	dueHour = 23
)

// ToCards returns a slice of trello cards from the given habits which haven't been marked today
func ToCards(habits map[string]habit.Habit, now time.Time) (cards []trello.Card, err error) {
	for name, habit := range habits {
		if habit.State != "" {
			continue
		}

		// include the day of month in card title to force overwrite in the beginning of the next day
		title := fmt.Sprintf("%v (%d)", name, now.Day())

		due := time.Date(now.Year(), now.Month(), now.Day(), dueHour, 0, 0, 0, now.Location())
		c, err := trello.NewCard(title, habit.CellName, &due)
		if err != nil {
			return nil, fmt.Errorf("could not create habit card: %w", err)
		}

		cards = append(cards, c)
	}

	return cards, nil
}
