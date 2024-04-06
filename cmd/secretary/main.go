package main

import (
	"database/sql"
	"errors"
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

	dbDriverName string = "postgres"
	envDBHost    string = "DB_HOST"
	envDBPort    string = "DB_PORT"
	envDBSSLMode string = "DB_SSLMODE"

	envDBName     string = "POSTGRES_DB"
	envDBUser     string = "POSTGRES_USER"
	envDBPassword string = "POSTGRES_PASSWORD"
)

type bot_app struct {
	bot *tgbotapi.BotAPI
	db  *sql.DB
}

type product struct {
	name       string
	sum        int
	paymentDay int
}

func helpMessage() string {
	return fmt.Sprintf(`Here's what I can do:
/%s to start application,
/%s to see help message,
/%s <name> <sum> <payment day> to add,
/%s to see how many dollars you spent on your next sallary.`, cmdStart, cmdHelp, cmdAddProduct, cmdGetTotal)
}

func newApp() *bot_app {
	return &bot_app{bot: nil, db: nil}
}

func newProduct(args []string) (*product, error) {
	if len(args) != 3 {
		return nil, errors.New("Not enough arguments!")
	}
	sum, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Second argument should be a number!")
	}
	day, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Third argument should be a number!")
	}
	return &product{name: args[0], sum: sum, paymentDay: day}, nil
}

func (p *product) string() string {
	return fmt.Sprintf("Name: %s,\nSum: %d,\nPayment day: %d", p.name, p.sum, p.paymentDay)
}

func (app *bot_app) init() (err error) {
	log.Print("Initializing the application...")
	if err := dotenv.Load(); err != nil {
		return errors.New(
			fmt.Sprintf("Maybe %s file not found. Err message: %s", envFileName, err.Error()),
		)
	}
	app.bot, err = tgbotapi.NewBotAPI(os.Getenv(envTelegramToken))
	if err != nil {
		return
	}
	app.db, err = sql.Open(
		dbDriverName,
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv(envDBHost),
			os.Getenv(envDBPort),
			os.Getenv(envDBUser),
			os.Getenv(envDBPassword),
			os.Getenv(envDBName),
			os.Getenv(envDBSSLMode)),
	)
	if err != nil {
		return err
	}
	if err = app.db.Ping(); err != nil {
		return err
	}
	log.Print("Database connected")
	log.Print("Application is initialized!")

	return nil
}

func (app *bot_app) close() {
	app.db.Close()
	log.Print("Application is closed!")
}

func (app *bot_app) sendMessage(msg string, chatId int64) {
	tg_msg := tgbotapi.NewMessage(chatId, msg)
	app.bot.Send(tg_msg)
}

func (app *bot_app) getUpdates() (tgbotapi.UpdatesChannel, error) {
	upd := tgbotapi.NewUpdate(updateOffset)
	upd.Timeout = timeout

	updates, err := app.bot.GetUpdatesChan(upd)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (app *bot_app) addProduct(p *product) error {
	insertProductStr := `INSERT INTO products(Name, Sum, PaymentDay) VALUES($1, $2, $3)`
	_, err := app.db.Exec(insertProductStr, p.name, p.sum, p.paymentDay)
	if err != nil {
		return err
	}
	return nil
}

func parseCommandArguments(args string) []string {
	return strings.Split(args, " ")
}

func (app *bot_app) getTotal() (int, error) {
	dayNow := time.Now().Day()
	var getTotalStr string
	if dayNow < 5 || dayNow >= 20 {
		getTotalStr = `SELECT SUM(sum) AS total FROM products WHERE paymentDay >= 5 AND paymentDay < 20`
	} else {
		getTotalStr = `SELECT SUM(sum) AS total FROM products WHERE paymentDay < 5 OR paymentDay >= 20`
	}
	rows, err := app.db.Query(getTotalStr)
	if err != nil {
		return 0, err
	}
	var total int
	rows.Next()
	rows.Scan(&total)

	return total, nil
}

func (app *bot_app) processAnUpdate(upd tgbotapi.Update) error {
	if upd.Message == nil {
		return errors.New("Not messages are consumed!")
	}
	switch upd.Message.Command() {
	case cmdStart:
		app.sendMessage(fmt.Sprintf("Application started! Try /%s", cmdHelp), upd.Message.Chat.ID)
	case cmdHelp:
		app.sendMessage(helpMessage(), upd.Message.Chat.ID)
	case cmdAddProduct:
		p, val_err := newProduct(parseCommandArguments(upd.Message.CommandArguments()))
		if val_err != nil {
			app.sendMessage(val_err.Error(), upd.Message.Chat.ID)
			return nil
		}
		if err := app.addProduct(p); err != nil {
			app.sendMessage(
				"Can't add product: Something went wrong on the server :(",
				upd.Message.Chat.ID,
			)
			return err
		}
		app.sendMessage(
			fmt.Sprintf("Successfuly added new product:\n\n%s", p.string()),
			upd.Message.Chat.ID,
		)
	case cmdGetTotal:
		total, err := app.getTotal()
		if err != nil {
			return err
		}
		app.sendMessage(fmt.Sprintf("Total: %d", total), upd.Message.Chat.ID)
	default:
		app.sendMessage(
			fmt.Sprintf("Command not recognized. Try /%s", cmdHelp),
			upd.Message.Chat.ID,
		)
	}
	return nil
}

func main() {
	app := newApp()
	if err := app.init(); err != nil {
		log.Fatal(err)
	}
	defer app.close()
	upds, err := app.getUpdates()
	if err != nil {
		log.Fatal(err)
	}
	for upd := range upds {
		if err := app.processAnUpdate(upd); err != nil {
			log.Print(err)
		}
	}
}
