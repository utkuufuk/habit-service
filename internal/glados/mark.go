package glados

import (
	"fmt"
	"regexp"

	"github.com/utkuufuk/habit-service/internal/habit"
)

func markHabit(client habit.Client, cell, symbol string) error {
	matched, err := regexp.MatchString(`[a-zA-Z]{3}\ 202\d\![A-Z][1-9][0-9]?$|^100$`, cell)
	if err != nil || matched == false {
		return fmt.Errorf("invalid cell '%s' to mark habit: %v", cell, err)
	}

	if !habit.IsValidMarkSymbol(symbol) {
		return fmt.Errorf("invalid habit symbol '%s' to mark habit", symbol)
	}

	if err := client.MarkHabit(cell, symbol); err != nil {
		return fmt.Errorf("could not mark habit on cell '%s' with symbol '%s': %v", cell, symbol, err)
	}

	return nil
}
