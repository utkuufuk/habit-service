package entrello

import (
	"testing"
	"time"

	"github.com/utkuufuk/habit-service/internal/habit"
)

func TestToCards(t *testing.T) {
	cellName := "Jun 2020!C3"

	tt := []struct {
		name     string
		habits   map[string]habit.Habit
		numCards int
		err      error
	}{
		{
			name: "marked habits",
			habits: map[string]habit.Habit{
				"a": {CellName: cellName, State: "✔", Score: 0},
				"b": {CellName: cellName, State: "x", Score: 0},
				"c": {CellName: cellName, State: "✘", Score: 0},
				"d": {CellName: cellName, State: "–", Score: 0},
				"e": {CellName: cellName, State: "-", Score: 0},
			},
			numCards: 0,
			err:      nil,
		},
		{
			name: "some marked some unmarked habits",
			habits: map[string]habit.Habit{
				"a": {CellName: cellName, State: "✔", Score: 0},
				"b": {CellName: cellName, State: "x", Score: 0},
				"c": {CellName: cellName, State: "✘", Score: 0},
				"d": {CellName: cellName, State: "–", Score: 0},
				"e": {CellName: cellName, State: "-", Score: 0},
				"f": {CellName: cellName, State: "", Score: 0},
				"g": {CellName: cellName, State: "", Score: 0},
			},
			numCards: 2,
			err:      nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := toCards(tc.habits, time.Now())
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if len(cards) != tc.numCards {
				t.Errorf("expected %d cards, got %d", tc.numCards, len(cards))
			}
		})
	}
}
