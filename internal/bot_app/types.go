package bot_app

import (
	"database/sql"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type bot_app struct {
	bot *tgbotapi.BotAPI
	db  *sql.DB
}

type expenseReportRow struct {
	Name string
	Sum  int
}

type expenseReport struct {
	rows  []expenseReportRow
	total expenseReportRow
}
