package main

import (
	"log"
	"os"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	log.Print("Initializing the application...")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Print(err)
		panic(err)
	}

	log.Print("Application is initialized!")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, i`m youre secretary.")
			bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here is how I can help you.")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I received your message, but didn't understand. Try /help")
			bot.Send(msg)
		}
	}

	log.Print("Application is closed!")
}
