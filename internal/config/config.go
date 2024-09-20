package config

import "time"

const (
	CmdStart         string = "start"
	CmdHelp          string = "help"
	CmdAdd           string = "add"
	CmdExpenseReport string = "report"
)

const (
	Timeout          int           = 60
	UpdateOffset     int           = 0
	PingDuration     time.Duration = 60
	HeartBitDuration time.Duration = 10
)

const (
	EnvFileName      string = ".env"
	EnvTelegramToken string = "TELEGRAM_TOKEN"

	DbDriverName string = "postgres"
	EnvDBHost    string = "DB_HOST"
	EnvDBPort    string = "DB_PORT"
	EnvDBSSLMode string = "DB_SSLMODE"

	EnvDBName     string = "POSTGRES_DB"
	EnvDBUser     string = "POSTGRES_USER"
	EnvDBPassword string = "POSTGRES_PASSWORD"
)
