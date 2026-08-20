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
