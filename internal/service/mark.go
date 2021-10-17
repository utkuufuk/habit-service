package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/utkuufuk/habit-service/internal/habit"
)

type MarkHabitAction struct {
	Cell   string
	Symbol string
}

func (a MarkHabitAction) Run(ctx context.Context, client habit.Client) (string, error) {
	matched, err := regexp.MatchString(`[a-zA-Z]{3}\ 202\d\![A-Z][1-9][0-9]?$|^100$`, a.Cell)
	if err != nil || matched == false {
		return "", fmt.Errorf("invalid cell '%s' to mark habit: %v", a.Cell, err)
	}

	if !habit.IsValidMarkSymbol(a.Symbol) {
		return "", fmt.Errorf("invalid habit symbol '%s' to mark habit", a.Symbol)
	}

	if err := client.MarkHabit(a.Cell, a.Symbol); err != nil {
		return "", fmt.Errorf("could not mark habit on cell '%s' with symbol '%s': %v", a.Cell, a.Symbol, err)
	}

	return "", nil
}
