package habit

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestCreateHabitMap(t *testing.T) {
	any := "."

	tt := []struct {
		name string
		rows [][]string
		out  map[string]Habit
		err  error
	}{
		{
			name: "blank name",
			rows: [][]string{
				{"", ""},
				{"–", "–"},
				{any, any},
				{any, "✔"},
			},
			out: nil,
			err: errors.New(""),
		},
		{
			name: "all marked",
			rows: [][]string{
				{"", "a", "b", "c"},
				{"", "40", "50", "60"},
				{any, "✔", "✘", "–"},
				{any, "✔", "✔", "✘"},
			},
			out: map[string]Habit{
				"a": {"Jan 2020!B4", "✔", 1},
				"b": {"Jan 2020!C4", "✔", 0.5},
				"c": {"Jan 2020!D4", "✘", 0},
			},
		},
		{
			name: "blank mid rows",
			rows: [][]string{
				{"", "a", "b", "c"},
				{"–", "–", "–", "–"},
				{any, "✔", "✘", "–"},
				{any, "✔", "✔", "✘"},
			},
			out: map[string]Habit{
				"a": {"Jan 2020!B4", "✔", 1},
				"b": {"Jan 2020!C4", "✔", 0.5},
				"c": {"Jan 2020!D4", "✘", 0},
			},
		},
		{
			name: "blank cell in the middle",
			rows: [][]string{
				{"", "a", "b", "c", "d"},
				{"–", "–", "–", "–", "–"},
				{any, "✔", "✘", "", ""},
				{any, "✔", "✘", "", "–"},
			},
			out: map[string]Habit{
				"a": {"Jan 2020!B4", "✔", 1},
				"b": {"Jan 2020!C4", "✘", 0},
				"c": {"Jan 2020!D4", "", 0},
				"d": {"Jan 2020!E4", "–", 0},
			},
		},
		{
			name: "blank cells in the end",
			rows: [][]string{
				{any, "a", "b", "c", "d", "e"},
				{"–", "–", "–", "–", "–", "–"},
				{any, "✔", "✘", ""},
				{any, "✔", "✘", "–"},
			},
			out: map[string]Habit{
				"a": {"Jan 2020!B4", "✔", 1},
				"b": {"Jan 2020!C4", "✘", 0},
				"c": {"Jan 2020!D4", "–", 0},
				"d": {"Jan 2020!E4", "", 0},
				"e": {"Jan 2020!F4", "", 0},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			date := time.Date(2020, time.Month(1), 2, 0, 0, 0, 0, time.UTC)

			data := make([][]interface{}, 0, len(tc.rows))
			for r, row := range tc.rows {
				data = append(data, make([]interface{}, 0, len(row)))
				for _, col := range row {
					data[r] = append(data[r], col)
				}
			}

			habits, err := createHabitMap(data, date)
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if diff := cmp.Diff(habits, tc.out); diff != "" {
				t.Errorf("output diff: %s", diff)
			}
		})
	}
}
