package main

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
				"a": {"a", cellName, "✔", 0},
				"b": {"b", cellName, "x", 0},
				"c": {"c", cellName, "✘", 0},
				"d": {"d", cellName, "–", 0},
				"e": {"e", cellName, "-", 0},
			},
			numCards: 0,
			err:      nil,
		},
		{
			name: "some marked some unmarked habits",
			habits: map[string]habit.Habit{
				"a": {"a", cellName, "✔", 0},
				"b": {"b", cellName, "x", 0},
				"c": {"c", cellName, "✘", 0},
				"d": {"d", cellName, "–", 0},
				"e": {"e", cellName, "-", 0},
				"f": {"f", cellName, "", 0},
				"g": {"g", cellName, "", 0},
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
