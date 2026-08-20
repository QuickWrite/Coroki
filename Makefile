APP_NAME := coroki
CMD_PATH := ./cmd/coroki
BIN_DIR := ./bin

.PHONY: build run test clean

build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test ./...

clean:
	rm -rf $(BIN_DIR)


# Development Postgres Database
.PHONY: db-up db-down db-logs

db-up:
	docker compose -f postgres.yaml up -d

db-down:
	docker compose -f postgres.yaml down

db-logs:
	docker compose -f postgres.yaml logs -f
