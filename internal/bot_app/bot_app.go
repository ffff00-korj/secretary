package bot_app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	tg "github.com/Syfaro/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/ffff00-korj/secretary/internal/config"
	"github.com/ffff00-korj/secretary/internal/product"
	"github.com/ffff00-korj/secretary/internal/utils"
)

type bot_app struct {
	bot *tg.BotAPI
	db  *sqlx.DB
}

type expenseReportRow struct {
	Name       string
	Sum        int
	PaymentDay int
}

type expenseReport struct {
	rows  []expenseReportRow
	total expenseReportRow
}

func NewApp() *bot_app {
	return &bot_app{bot: nil, db: nil}
}

func (app *bot_app) Init() (err error) {
	log.Print("Initializing the application...")
	if err := dotenv.Load(); err != nil {
		return errors.New(
			fmt.Sprintf(
				"Maybe %s file not found: %s",
				config.EnvFileName,
				err.Error(),
			),
		)
	}
	app.bot, err = tg.NewBotAPI(os.Getenv(config.EnvTelegramToken))
	if err != nil {
		return
	}
	app.db, err = sqlx.Connect(
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
	runPingDBTask(app.db, config.PingDuration, config.HeartBitDuration, config.HeartBitAttempts)
	log.Print("Database connected")
	log.Print("Application is initialized!")

	return nil
}

func runPingDBTask(
	db *sqlx.DB,
	pingDuration, heartBitDuration time.Duration,
	heartBitAttempts int,
) {
	ticker := time.NewTicker(pingDuration * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := db.Ping()
				if err == nil {
					continue
				} else {
					log.Printf("%s. Starting heart bits.", err)
				}
				heartBitTicker := time.NewTicker(pingDuration * time.Second)
				for tickCount := 1; tickCount <= heartBitAttempts; tickCount++ {
					select {
					case <-ticker.C:
						err = db.Ping()
						if err == nil {
							log.Println("DB connection restored.")
							break
						}
						log.Printf("%s. Heart bit %d", err, tickCount)
					case <-quit:
						heartBitTicker.Stop()
						return
					}
				}
				if err != nil {
					log.Panic("DB is not responding.")
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (app *bot_app) Close() {
	app.db.Close()
	log.Print("Application is closed!")
}

func (app *bot_app) GetUpdates() (tg.UpdatesChannel, error) {
	upd := tg.NewUpdate(config.UpdateOffset)
	upd.Timeout = config.Timeout

	updates, err := app.bot.GetUpdatesChan(upd)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (app *bot_app) ProcessAnUpdate(ctx context.Context, upd tg.Update) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
		p, err := product.NewProductFromArgs(upd.Message.CommandArguments())
		if err != nil {
			app.sendMessage(
				err.Error(),
				upd.Message.Chat.ID,
				"",
			)
			return fmt.Errorf(
				"Validation error: %w. args [%s]",
				err,
				upd.Message.CommandArguments(),
			)
		}
		exists, err := app.checkProductExists(p)
		if err != nil {
			app.sendMessage(
				"Can't add product: Something went wrong on the server :(",
				upd.Message.Chat.ID,
				"",
			)
			return err
		}
		if exists {
			app.sendMessage("Product with this name already exists.", upd.Message.Chat.ID, "")
			return nil
		}
		if err := app.addProduct(p); err != nil {
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
			app.sendMessage(
				"Can't get expense report: Something went wrong on the server :(",
				upd.Message.Chat.ID,
				"",
			)
			return err
		}
		app.sendMessage(utils.TextToMarkdown(report), upd.Message.Chat.ID, "markdown")
	default:
		app.sendMessage(
			fmt.Sprintf("Command not recognized. Try /%s", config.CmdHelp),
			upd.Message.Chat.ID,
			"",
		)
	}

	return nil
}

func (app *bot_app) ProcessAnUpdateWithContext(ctx context.Context, upds tg.Update) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errC := make(chan error)
	go func() {
		if err := app.ProcessAnUpdate(ctx, upds); err != nil {
			errC <- err
		}
	}()

	select {
	case <-ctx.Done():
        if ctx.Err() != nil {
		    log.Print(ctx.Err())
        }
	case err := <-errC:
		log.Print(err)
	}
}

func (app *bot_app) sendMessage(msg string, chatId int64, parser string) {
	tg_msg := tg.NewMessage(chatId, msg)
	if parser != "" {
		tg_msg.ParseMode = parser
	}
	app.bot.Send(tg_msg)
}
