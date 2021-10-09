package habit

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"google.golang.org/api/sheets/v4"
)

const (
	nameRowIdx    = 0 // number of rows before the name row starts in the spreadsheet
	scoreRowIdx   = 1 // number of rows before the score row starts in the spreadsheet
	dataRowIdx    = 2 // number of rows before the first data row starts in the spreadsheet
	dataColumnIdx = 1 // number of columns before the first data column starts in the spreadsheet

	symbolDone    = "✔"
	symbolFailed  = "✘"
	symbolSkipped = "–"
)

type Client struct {
	spreadsheetId string
	service       *sheets.SpreadsheetsValuesService
	location      *time.Location
}

type Habit struct {
	Name     string
	CellName string
	State    string
	Score    float64
}

type cell struct {
	col string
	row int
}

func GetClient(
	ctx context.Context,
	spreadsheetId string,
	location *time.Location,
) (client Client, err error) {
	service, err := initializeService(ctx)
	if err != nil {
		return client, fmt.Errorf("could not initialize gsheets service: %w", err)
	}
	return Client{spreadsheetId, service.Spreadsheets.Values, location}, nil
}

// FetchHabits retrieves the state of today's habits from the spreadsheet
func (c Client) FetchHabits(now time.Time) (map[string]Habit, error) {
	rangeName, err := getRangeName(now, cell{"A", 1}, cell{"Z", now.Day() + dataRowIdx})
	if err != nil {
		return nil, fmt.Errorf("could not get range name: %w", err)
	}

	rows, err := c.readCells(rangeName)
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}

	return mapHabits(rows, now)
}

func (c Client) MarkHabit(cellName string, symbol string) error {
	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, 1)
	values[0][0] = symbol
	return c.writeCells(values, cellName)
}

func IsValidMarkSymbol(symbol string) bool {
	return symbol == symbolDone || symbol == symbolFailed || symbol == symbolSkipped
}

func (c Client) UpdateScores(habits map[string]Habit, now time.Time) error {
	// map habit scores into a slice ordered by column
	scores := make([]float64, len(habits))
	var cellNameComponents []string
	for _, habit := range habits {
		cellNameComponents = strings.Split(habit.CellName, "!")
		col := []rune(cellNameComponents[1][0:1])[0]
		idx := int(col) - int('A') - 1
		scores[idx] = habit.Score
	}

	// get range name to write habit scores in the sheet
	row := scoreRowIdx + 1
	firstCol := string(rune(int('A') + dataColumnIdx))
	lastCol := string(rune(int('A') + len(habits)))
	rangeName, err := getRangeName(now, cell{firstCol, row}, cell{lastCol, row})
	if err != nil {
		return fmt.Errorf("could not get range name: %w", err)
	}

	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, len(scores))
	for i, score := range scores {
		values[0][i] = score
	}
	return c.writeCells(values, rangeName)
}

// readCells reads a range of cell values with the given range
func (c Client) readCells(rangeName string) ([][]interface{}, error) {
	resp, err := c.service.Get(c.spreadsheetId, rangeName).Do()
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}
	return resp.Values, nil
}

// writeCells writes a 2D array of values into a range of cells
func (c Client) writeCells(values [][]interface{}, rangeName string) error {
	_, err := c.service.
		Update(c.spreadsheetId, rangeName, &sheets.ValueRange{Values: values}).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return fmt.Errorf("could not write cells: %w", err)
	}
	return nil
}

// mapHabits creates a map of habits for given a date and a spreadsheet row data
func mapHabits(rows [][]interface{}, date time.Time) (map[string]Habit, error) {
	habits := make(map[string]Habit)
	for col := dataColumnIdx; col < len(rows[0]); col++ {
		// evaluate the habit's cell name for today
		c := cell{string(rune('A' + col)), date.Day() + dataRowIdx}
		cellName, err := getRangeName(date, c, c)
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
			if val == symbolDone {
				nom++
			}

			if val == symbolSkipped {
				denom--
			}
		}
		score := (float64(nom) / float64(denom))
		if math.IsNaN(score) {
			score = 0
		}

		habits[name] = Habit{name, cellName, state, score}
	}
	return habits, nil
}

// getRangeName gets the range name given a date and start & end cells
func getRangeName(date time.Time, start, end cell) (string, error) {
	if start.col < "A" || start.col > "Z" || start.row <= 0 {
		return "", fmt.Errorf("invalid start cell: %s%d", start.col, start.row)
	}

	month := date.Month().String()[:3]
	year := date.Year()

	// assume single cell if no end date specified
	if end.col == "" || end.row == 0 || (end.col == start.col && end.row == start.row) {
		return fmt.Sprintf("%s %d!%s%d", month, year, start.col, start.row), nil
	}

	if end.col < "A" || end.col > "Z" || end.row <= 0 {
		return "", fmt.Errorf("invalid end cell: %s%d", end.col, end.row)
	}

	return fmt.Sprintf("%s %d!%s%d:%s%d", month, year, start.col, start.row, end.col, end.row), nil
}
