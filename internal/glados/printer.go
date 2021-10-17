package glados

import (
	"strconv"

	"github.com/utkuufuk/habit-service/internal/tableimage"
)

type table struct {
	rows []tableimage.TR
}

func newTable() *table {
	return &table{make([]tableimage.TR, 0)}
}

func (t *table) save(path string) {
	ti := tableimage.Init("#fff", tableimage.Png, path)

	ti.AddTH(
		tableimage.TR{
			BorderColor: "#000",
			Tds: []tableimage.TD{
				{
					Color: "#000",
					Text:  "Habit",
				},
				{
					Color: "#000",
					Text:  "Last Month",
				},
				{
					Color: "#000",
					Text:  "This Month",
				},
				{
					Color: "#000",
					Text:  "Delta",
				},
			},
		},
	)

	ti.AddTRs(t.rows)
	ti.Save()
}

func (t *table) addRow(name string, last, current float64) {
	deltaColor := "#C82538"
	if current > last {
		deltaColor = "#2E7F18"
	}

	row := tableimage.TR{
		BorderColor: "#000",
		Tds: []tableimage.TD{
			{
				Color: "#000",
				Text:  name,
			},
			{
				Color: getScoreColor(last),
				Text:  strconv.FormatFloat(last, 'f', 0, 32) + "%",
			},
			{
				Color: getScoreColor(current),
				Text:  strconv.FormatFloat(current, 'f', 0, 32) + "%",
			},
			{
				Color: deltaColor,
				Text:  strconv.FormatFloat(current-last, 'f', 0, 32) + "%",
			},
		},
	}
	t.rows = append(t.rows, row)
}

func getScoreColor(score float64) string {
	if score > 83 {
		return "#2E7F18"
	}

	if score > 67 {
		return "#45731E"
	}

	if score > 50 {
		return "#675E24"
	}

	if score > 33 {
		return "#8D472B"
	}

	if score > 16 {
		return "#B13433"
	}

	return "#C82538"
}
