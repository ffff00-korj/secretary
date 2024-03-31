package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	cmdStart      string = "start"
	cmdHelp       string = "help"
	cmdAddProduct string = "add"
	cmdGetTotal   string = "total"
)

const (
	timeout      int = 60
	updateOffset int = 0
)

const (
	envFileName      string = ".env"
	envTelegramToken string = "TELEGRAM_TOKEN"

	envDBHost    string = "DB_HOST"
	envDBPort    string = "DB_PORT"
	envDBSSLMode string = "DB_SSLMODE"

	envDBName     string = "POSTGRES_DB"
	envDBUser     string = "POSTGRES_USER"
	envDBPassword string = "POSTGRES_PASSWORD"
)

type product struct {
	name       string
	sum        int
	paymentDay int
}

type products []product

func app() (*tgbotapi.BotAPI, *sql.DB) {
	if err := dotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("No %s file found", envFileName))
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv(envTelegramToken))
	if err != nil {
		log.Fatal(err)
	}
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv(envDBHost),
		os.Getenv(envDBPort),
		os.Getenv(envDBUser),
		os.Getenv(envDBPassword),
		os.Getenv(envDBName),
		os.Getenv(envDBSSLMode),
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Database connected")

	return bot, db
}

func getUpdate(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(updateOffset)
	u.Timeout = timeout

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	return updates
}

func addProduct(db *sql.DB, p product) {
	insertProductStr := `INSERT INTO products(Name, Sum, PaymentDay) VALUES($1, $2, $3)`
	_, err := db.Exec(insertProductStr, p.name, p.sum, p.paymentDay)
	if err != nil {
		log.Fatal(err)
	}
}

func parseCommandArguments(args string) []string {
	return strings.Split(args, " ")
}

func getTotal(db *sql.DB) int {
	dayNow := time.Now().Day()
	var getTotalStr string
	if dayNow < 5 || dayNow >= 20 {
		getTotalStr = `SELECT SUM(sum) AS total FROM products WHERE paymentDay >= 5 AND paymentDay < 20`
	} else {
		getTotalStr = `SELECT SUM(sum) AS total FROM products WHERE paymentDay < 5 OR paymentDay >= 20`
	}
	rows, err := db.Query(getTotalStr)
	if err != nil {
		log.Fatal(err)
	}
	var total int
	rows.Next()
	rows.Scan(&total)

	return total
}

func processAnUpdate(bot *tgbotapi.BotAPI, db *sql.DB, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case cmdStart:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf(
					"Application started! Try /%s",
					cmdHelp,
				),
			)
			bot.Send(msg)
		case cmdHelp:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				fmt.Sprintf(`Here's what I can do:
/%s to start application,
/%s to see help message,
/%s <name> <sum> <payment day> to add,
/%s to see how many dollars you spent on your next sallary.`,
					cmdStart, cmdHelp, cmdAddProduct, cmdGetTotal))
			bot.Send(msg)
		case cmdAddProduct:
			arguments := parseCommandArguments(update.Message.CommandArguments())
			if len(arguments) != 3 {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Not enough arguments!",
				)
				bot.Send(msg)
				continue
			}
			sum, err := strconv.Atoi(arguments[1])
			if err != nil {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Second argument should be a number!",
				)
				bot.Send(msg)
				continue
			}
			day, err := strconv.Atoi(arguments[2])
			if err != nil {
				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"Third argument should be a number!",
				)
				bot.Send(msg)
				continue
			}
			addProduct(db, product{name: arguments[0], sum: sum, paymentDay: day})
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Successfuly added new Product!",
			)
			bot.Send(msg)
		case cmdGetTotal:
			total := getTotal(db)
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Total: %d", total),
			)
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Command not recognized. Try /%s", cmdHelp),
			)
			bot.Send(msg)
		}
	}
}

func main() {
	log.Print("Initializing the application...")
	bot, db := app()
	defer db.Close()
	log.Print("Application is initialized!")
	updates := getUpdate(bot)
	processAnUpdate(bot, db, updates)
	log.Print("Application is closed!")
}
