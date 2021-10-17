package glados

import (
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

func (t *table) addRow(name, last, current, delta string) {
	row := tableimage.TR{
		BorderColor: "#000",
		Tds: []tableimage.TD{
			{
				Color: "#000",
				Text:  name,
			},
			{
				Color: "#000",
				Text:  last + "%",
			},
			{
				Color: "#000",
				Text:  current + "%",
			},
			{
				Color: "#000",
				Text:  delta + "%",
			},
		},
	}
	t.rows = append(t.rows, row)
}
