package habit

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/utkuufuk/habit-service/internal/sheets"
)

const (
	nameRowIdx    = 0 // number of rows before the name row starts in the spreadsheet
	scoreRowIdx   = 1 // number of rows before the score row starts in the spreadsheet
	dataRowIdx    = 2 // number of rows before the first data row starts in the spreadsheet
	dataColumnIdx = 1 // number of columns before the first data column starts in the spreadsheet

	SymbolDone = "✔"
	SymbolFail = "✘"
	SymbolSkip = "–"
)

type Habit struct {
	Cell  sheets.Range
	State string
	Score float64
}

// FetchAll retrieves the state of habits from the spreadsheet as of the selected date
func FetchAll(client sheets.Client, date time.Time) (map[string]Habit, error) {
	rng, err := sheets.GetRange(
		getSheetName(date),
		sheets.Cell{Col: "A", Row: 1},
		sheets.Cell{Col: "Z", Row: date.Day() + dataRowIdx},
	)
	if err != nil {
		return nil, fmt.Errorf("could not get spreadsheet range: %w", err)
	}

	rows, err := client.ReadCells(rng)
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}

	return parseHabitMap(rows, date)
}

// Mark marks the habit with the given cell as done/failed/skipped
func Mark(client sheets.Client, cell, symbol string) error {
	ok, err := regexp.MatchString(`[a-zA-Z]{3} 202\d![A-Z]([1-9]|[1-3][0-9])$`, cell)
	if err != nil || ok == false {
		return fmt.Errorf("invalid cell name '%s': %w", cell, err)
	}

	if symbol != SymbolDone && symbol != SymbolFail && symbol != SymbolSkip {
		return fmt.Errorf("invalid symbol '%s' to mark habit", symbol)
	}

	if err := client.WriteCell(cell, symbol); err != nil {
		return fmt.Errorf("could not write symbol '%s' to cell '%s': %w", symbol, cell, err)
	}

	return nil
}

// WriteScores (re)writes the scores of the given habits in the spreadsheet
func WriteScores(client sheets.Client, date time.Time, habits []Habit) error {
	scores := make([]float64, len(habits))
	for _, habit := range habits {
		idx := habit.Cell.GetStartColumnIndex()
		scores[idx] = habit.Score
	}

	row := scoreRowIdx + 1
	scoreRowRange, err := sheets.GetRange(
		getSheetName(date),
		sheets.Cell{Col: string(rune(int('A') + dataColumnIdx)), Row: row},
		sheets.Cell{Col: string(rune(int('A') + len(habits))), Row: row},
	)
	if err != nil {
		return fmt.Errorf("could not get score row range: %w", err)
	}

	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, len(scores))
	for i, score := range scores {
		values[0][i] = score
	}
	return client.WriteCells(values, scoreRowRange)
}

// getSheetName gets the sheet name corresponding to the given date
func getSheetName(date time.Time) string {
	month := date.Month().String()[:3]
	year := date.Year()
	return fmt.Sprintf("%s %d", month, year)
}

// parseHabitMap parses a map of habits from spreadsheet row data for the given date
func parseHabitMap(rows [][]interface{}, date time.Time) (map[string]Habit, error) {
	habits := make(map[string]Habit)
	for col := dataColumnIdx; col < len(rows[0]); col++ {
		// evaluate the habit's cell name the selected date
		c := sheets.Cell{Col: string(rune('A' + col)), Row: date.Day() + dataRowIdx}
		cell, err := sheets.GetRange(getSheetName(date), c, c)
		if err != nil {
			return nil, err
		}

		// handle cases where the last N columns are blank which reduces the slice length by N
		state := ""
		if col < len(rows[date.Day()+dataRowIdx-1]) {
			state = fmt.Sprintf("%v", rows[date.Day()+dataRowIdx-1][col])
		}

		// read habit name
		name := fmt.Sprintf("%v", rows[nameRowIdx][col])
		if name == "" {
			return nil, fmt.Errorf("habit name cannot be blank")
		}

		// calculate habit score
		nom := 0
		denom := date.Day()
		for row := dataRowIdx; row < date.Day()+dataRowIdx; row++ {
			if len(rows[row]) < col+1 {
				continue
			}

			val := rows[row][col]
			if val == SymbolDone {
				nom++
			}

			if val == SymbolSkip {
				denom--
			}
		}
		score := (float64(nom) / float64(denom))
		if math.IsNaN(score) {
			score = 0
		}

		habits[name] = Habit{cell, state, score}
	}
	return habits, nil
}
