package service

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/pkg/trello"
	"github.com/utkuufuk/habit-service/internal/entrello"
	"github.com/utkuufuk/habit-service/internal/habit"
)

type FetchHabitsAsTrelloCardsAction struct {
	TimezoneLocation *time.Location
}

func (a FetchHabitsAsTrelloCardsAction) Run(ctx context.Context, client habit.Client) ([]trello.Card, error) {
	now := time.Now().In(a.TimezoneLocation)

	habits, err := client.FetchHabits(now)
	if err != nil {
		return nil, fmt.Errorf("could not fetch habits: %w", err)
	}

	// @todo: call this as part of a cron job
	if err = client.UpdateScores(habits, now); err != nil {
		return nil, fmt.Errorf("could not update habit scores: %w", err)
	}

	return entrello.ToCards(habits, now)
}
