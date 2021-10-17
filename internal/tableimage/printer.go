package tableimage

import (
	"strconv"
)

type table struct {
	rows []TR
}

func NewTable() *table {
	return &table{make([]TR, 0)}
}

func (t *table) Save(path string) {
	ti := Init("#fff", Png, path)

	ti.AddTH(
		TR{
			BorderColor: "#000",
			Tds: []TD{
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

func (t *table) AddRow(name string, last, current float64) {
	deltaColor := "#C82538"
	if current > last {
		deltaColor = "#2E7F18"
	}

	row := TR{
		BorderColor: "#000",
		Tds: []TD{
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
