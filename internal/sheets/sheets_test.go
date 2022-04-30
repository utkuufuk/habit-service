package sheets

import (
	"errors"
	"testing"
)

func TestGetRangeName(t *testing.T) {
	testSheetName := "Jan 2020"

	tt := []struct {
		name      string
		sheetName string
		start     Cell
		end       Cell
		out       string
		err       error
	}{
		{
			name:      "invalid start col",
			sheetName: testSheetName,
			start:     Cell{"", 1},
			end:       Cell{},
			err:       errors.New(""),
		},
		{
			name:      "invalid start row",
			sheetName: testSheetName,
			start:     Cell{"A", 0},
			end:       Cell{},
			err:       errors.New(""),
		},
		{
			name:      "invalid end col",
			sheetName: testSheetName,
			start:     Cell{"A", 1},
			end:       Cell{"0", 1},
			err:       errors.New(""),
		},
		{
			name:      "invalid end row",
			sheetName: testSheetName,
			start:     Cell{"A", 1},
			end:       Cell{"A", -1},
			err:       errors.New(""),
		},
		{
			name:      "implicit single cell",
			sheetName: testSheetName,
			start:     Cell{"A", 1},
			end:       Cell{},
			out:       "Jan 2020!A1",
			err:       nil,
		},
		{
			name:      "explicit single cell",
			sheetName: testSheetName,
			start:     Cell{"A", 1},
			end:       Cell{"A", 1},
			out:       "Jan 2020!A1",
			err:       nil,
		},
		{
			name:      "valid range",
			sheetName: testSheetName,
			start:     Cell{"B", 3},
			end:       Cell{"D", 5},
			out:       "Jan 2020!B3:D5",
			err:       nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rng, err := GetRange(testSheetName, tc.start, tc.end)
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if err == nil {
				return
			}

			if string(rng) != tc.out {
				t.Fatalf("range name mismatch; want '%s', got '%v'", tc.out, rng)
			}
		})
	}
}
