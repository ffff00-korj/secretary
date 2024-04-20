include .env

PACKAGE_NAME = secretary
PACKAGE_PATH = ./cmd/$(PACKAGE_NAME)/main.go
BUILD_PATH = ./build
MIGRATIONS_DIR = ./migrations

build_delete:
	rm -rf build

build: build_delete
	go build -o $(BUILD_PATH)/$(PACKAGE_NAME) $(PACKAGE_PATH)

run:
	go run $(PACKAGE_PATH)

migrate_up:
	goose -dir $(MIGRATIONS_DIR) postgres "host=$(DB_HOST) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=$(DB_SSLMODE)" up

db_run:
	docker compose up -d

db_down:
	docker compose down

db_dump:
	pg_dump -U $(POSTGRES_USER) -h $(DB_HOST) -d $(POSTGRES_DB) > $(DB_BACKUP_PATH)/dump.sql
