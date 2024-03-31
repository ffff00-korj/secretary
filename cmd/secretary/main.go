package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	cmdStart      string = "start"
	cmdHelp       string = "help"
	cmdAddProduct string = "add"
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

func app() *tgbotapi.BotAPI {
	if err := dotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("No %s file found", envFileName))
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv(envTelegramToken))
	if err != nil {
		log.Fatal(err)
	}
	return bot
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

func addTestProduct(name string) {
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
	insertProductStr := `INSERT INTO products(Name) VALUES($1)`
	_, err = db.Exec(insertProductStr, name)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func processAnUpdate(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Command() {
		case cmdStart:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf(
					"Hi, i`m been employed as your personal assistant. If you want to know what i can do type %s",
					cmdHelp,
				),
			)
			bot.Send(msg)
		case cmdHelp:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here's what I can do")
			bot.Send(msg)
		case cmdAddProduct:
			addTestProduct(string(update.Message.CommandArguments()))
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Successfuly added new Product!",
			)
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("I received your message, but did not recognize it. Try %s", cmdHelp),
			)
			bot.Send(msg)
		}
	}
}

func main() {
	log.Print("Initializing the application...")
	bot := app()
	log.Print("Application is initialized!")
	updates := getUpdate(bot)
	processAnUpdate(bot, updates)
	log.Print("Application is closed!")
}
