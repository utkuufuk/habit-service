package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/sheets"
	"github.com/utkuufuk/habit-service/internal/tableimage"
	"golang.org/x/exp/slices"
)

type table struct {
	rows []tableimage.TR
}

func ReportProgress(
	client sheets.Client,
	loc *time.Location,
	skipList []string,
	telegramChatId int64,
	telegramToken string,
) error {
	now := time.Now().In(loc)
	thisMonthHabits, err := habit.FetchAll(client, now)
	if err != nil {
		return fmt.Errorf("could not fetch this month's habits: %w\n", err)
	}

	year, month, _ := now.Date()
	endOfLastMonth := time.Date(year, month, 1, 0, 0, 0, 0, loc).Add(-time.Nanosecond)
	lastMonthHabits, err := habit.FetchAll(client, endOfLastMonth)
	if err != nil {
		return fmt.Errorf("could not fetch habits from last month: %w\n", err)
	}

	table := &table{make([]tableimage.TR, 0)}
	for name, habit := range thisMonthHabits {
		if slices.Contains(skipList, name) {
			continue
		}

		table.addRow(name, lastMonthHabits[name].Score*100, habit.Score*100)
	}

	path := fmt.Sprintf("./progress-report-%s.png", now.Format("2006-01-02T15:04:05"))
	table.save(path)
	err = sendProgressReport(path, telegramChatId, telegramToken)
	if err != nil {
		return fmt.Errorf("could not send progress report: %w\n", err)
	}

	return os.Remove(path)
}

func sendProgressReport(path string, telegramChatId int64, telegramToken string) error {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		return fmt.Errorf("could not initialize Telegram bot client: %w", err)
	}

	photoBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read progress report image '%s': %w", path, err)
	}

	_, err = bot.Send(tgbotapi.NewPhotoUpload(telegramChatId, tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: photoBytes,
	}))
	return err
}

func (t *table) save(path string) {
	header := tableimage.TR{
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
	}
	tableimage.Draw(tableimage.TableImage{
		BackgroundColor: "#fff",
		FileType:        tableimage.Png,
		FilePath:        path,
		Header:          header,
		Rows:            t.rows,
	})
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
