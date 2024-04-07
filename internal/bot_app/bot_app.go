package bot_app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ffff00-korj/secretary/internal/config"
	"github.com/ffff00-korj/secretary/internal/product"
	"github.com/ffff00-korj/secretary/internal/utils"
)

func NewApp() *bot_app {
	return &bot_app{bot: nil, db: nil}
}

func (app *bot_app) Init() (err error) {
	log.Print("Initializing the application...")
	if err := dotenv.Load(); err != nil {
		return errors.New(
			fmt.Sprintf(
				"Maybe %s file not found. Err message: %s",
				config.EnvFileName,
				err.Error(),
			),
		)
	}
	app.bot, err = tgbotapi.NewBotAPI(os.Getenv(config.EnvTelegramToken))
	if err != nil {
		return
	}
	app.db, err = sql.Open(
		config.DbDriverName,
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv(config.EnvDBHost),
			os.Getenv(config.EnvDBPort),
			os.Getenv(config.EnvDBUser),
			os.Getenv(config.EnvDBPassword),
			os.Getenv(config.EnvDBName),
			os.Getenv(config.EnvDBSSLMode)),
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

func (app *bot_app) Close() {
	app.db.Close()
	log.Print("Application is closed!")
}

func (app *bot_app) GetUpdates() (tgbotapi.UpdatesChannel, error) {
	upd := tgbotapi.NewUpdate(config.UpdateOffset)
	upd.Timeout = config.Timeout

	updates, err := app.bot.GetUpdatesChan(upd)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (app *bot_app) ProcessAnUpdate(upd tgbotapi.Update) error {
	if upd.Message == nil {
		return errors.New("No messages are consumed!")
	}
	switch upd.Message.Command() {
	case config.CmdStart:
		app.sendMessage(
			fmt.Sprintf("Application started! Try /%s", config.CmdHelp),
			upd.Message.Chat.ID,
			"",
		)
	case config.CmdHelp:
		app.sendMessage(utils.HelpMessage(), upd.Message.Chat.ID, "")
	case config.CmdAdd:
		p, val_err := product.NewProduct(
			utils.ParseMessageArguments(upd.Message.CommandArguments()),
		)
		if val_err != nil {
			app.sendMessage(val_err.Error(), upd.Message.Chat.ID, "")
			return nil
		}
		exists, err := app.checkProductExists(p)
		if err != nil {
			return err
		}
		if exists {
			app.sendMessage("Product with this name already exists.", upd.Message.Chat.ID, "")
			return nil
		}
		if _, err := app.addProduct(p); err != nil {
			app.sendMessage(
				"Can't add product: Something went wrong on the server :(",
				upd.Message.Chat.ID,
				"",
			)
			return err
		}
		app.sendMessage(
			fmt.Sprintf("Successfuly added new product:\n\n%s", p.String()),
			upd.Message.Chat.ID,
			"",
		)
	case config.CmdExpenseReport:
		report, err := app.getExpenseReport()
		if err != nil {
			return err
		}
		app.sendMessage("```"+report+"```", upd.Message.Chat.ID, "markdown")
	default:
		app.sendMessage(
			fmt.Sprintf("Command not recognized. Try /%s", config.CmdHelp),
			upd.Message.Chat.ID,
			"",
		)
	}
	return nil
}

func (app *bot_app) sendMessage(msg string, chatId int64, parser string) {
	tg_msg := tgbotapi.NewMessage(chatId, msg)
	if parser != "" {
		tg_msg.ParseMode = parser
	}
	app.bot.Send(tg_msg)
}
