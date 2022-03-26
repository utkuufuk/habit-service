package service

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/entrello"
	"github.com/utkuufuk/habit-service/internal/habit"
)

type FetchHabitsAsTrelloCardsAction struct{}

func (a FetchHabitsAsTrelloCardsAction) Run(ctx context.Context, client habit.Client) ([]trello.Card, error) {
	now := time.Now().In(config.TimezoneLocation)

	habits, err := client.FetchHabits(now)
	if err != nil {
		return nil, fmt.Errorf("could not fetch habits: %w", err)
	}

	return entrello.ToCards(habits, now)
}
