package sheets

import (
	"fmt"
	"strings"

	"google.golang.org/api/sheets/v4"
)

type Cell struct {
	Col string
	Row int
}

type CellName string
type RangeName string

func (r CellName) GetStartColumnIndex() int {
	cellNameComponents := strings.Split(string(r), "!")
	col := []rune(cellNameComponents[1][0:1])[0]
	return int(col) - int('A') - 1
}

// GetCellName gets the cell name for the given cell
func GetCellName(sheetName string, cell Cell) (CellName, error) {
	rng, err := GetRangeName(sheetName, cell, cell)
	return CellName(rng), err
}

// GetRangeName gets the range name given the sheet name and start & end cells
func GetRangeName(sheetName string, start, end Cell) (rng RangeName, err error) {
	if start.Col < "A" || start.Col > "Z" || start.Row <= 0 {
		return rng, fmt.Errorf("invalid start cell: %s%d", start.Col, start.Row)
	}

	// return single cell name if no end cell specified
	if end.Col == "" || end.Row == 0 || (end.Col == start.Col && end.Row == start.Row) {
		return RangeName(fmt.Sprintf("%s!%s%d", sheetName, start.Col, start.Row)), nil
	}

	if end.Col < "A" || end.Col > "Z" || end.Row <= 0 {
		return rng, fmt.Errorf("invalid end cell: %s%d", end.Col, end.Row)
	}

	return RangeName(fmt.Sprintf("%s!%s%d:%s%d", sheetName, start.Col, start.Row, end.Col, end.Row)), nil
}

// ReadCells reads a range of cell values with the given range
func (c Client) ReadCells(rng RangeName) ([][]interface{}, error) {
	resp, err := c.service.Get(c.spreadsheetId, string(rng)).Do()
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}
	return resp.Values, nil
}

// WriteCell writes the given value into the cell
func (c Client) WriteCell(cellName string, value string) error {
	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, 1)
	values[0][0] = value
	return c.WriteCells(values, RangeName(cellName))
}

// writeCells writes a 2D array of values into a range of cells
func (c Client) WriteCells(values [][]interface{}, rng RangeName) error {
	_, err := c.service.
		Update(c.spreadsheetId, string(rng), &sheets.ValueRange{Values: values}).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return fmt.Errorf("could not write cells: %w", err)
	}
	return nil
}
