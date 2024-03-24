package main

import (
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/joho/godotenv"
)

const (
	timeout      int = 60
	updateOffset int = 0
)

func app() *tgbotapi.BotAPI {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Print(err)
		panic(err)
	}
	return bot
}

func getUpdate(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(updateOffset)
	u.Timeout = timeout

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	return updates
}

func processAnUpdate(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Hi, i`m been employed as your personal assistant. If you want to know what i can do type /help",
			)
			bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Here's what I can do")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"I received your message, but did not recognize it. Try /help",
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
