package service

import (
	"testing"
	"time"

	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/sheets"
)

func TestToTrelloCards(t *testing.T) {
	testCellName := sheets.CellName("Jun 2020!C3")

	tt := []struct {
		name     string
		habits   map[string]habit.Habit
		numCards int
		err      error
	}{
		{
			name: "marked habits",
			habits: map[string]habit.Habit{
				"a": {CellName: testCellName, State: "✔", Score: 0},
				"b": {CellName: testCellName, State: "x", Score: 0},
				"c": {CellName: testCellName, State: "✘", Score: 0},
				"d": {CellName: testCellName, State: "–", Score: 0},
				"e": {CellName: testCellName, State: "-", Score: 0},
			},
			numCards: 0,
			err:      nil,
		},
		{
			name: "some marked some unmarked habits",
			habits: map[string]habit.Habit{
				"a": {CellName: testCellName, State: "✔", Score: 0},
				"b": {CellName: testCellName, State: "x", Score: 0},
				"c": {CellName: testCellName, State: "✘", Score: 0},
				"d": {CellName: testCellName, State: "–", Score: 0},
				"e": {CellName: testCellName, State: "-", Score: 0},
				"f": {CellName: testCellName, State: "", Score: 0},
				"g": {CellName: testCellName, State: "", Score: 0},
			},
			numCards: 2,
			err:      nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := toTrelloCards(tc.habits, time.Now())
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if len(cards) != tc.numCards {
				t.Errorf("expected %d cards, got %d", tc.numCards, len(cards))
			}
		})
	}
}
