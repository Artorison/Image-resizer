MAIN := ./cmd/main.go
BIN_DIR := bin
TARGET := image-resizer
BIN   := $(BIN_DIR)/$(TARGET)
INTERNAL := ./internal/
PKG := ./pkg/

CONFIG_FILE := ./configs/config.yaml

.PHONY: all run build lint test clean test_integration up down check

all: run

run: build
	@$(BIN) --config $(CONFIG_FILE)

up:
	$(DC) up

down:
	$(DC) down

build:
	@go build -v -o $(BIN) $(MAIN)

lint:
	@golangci-lint run ./...

clean:
	rm -rf ./test/testdata
	rm -rf $(BIN_DIR)
	$(DC) down
	$(DC) -f $(DC_TEST) down

check: build lint test test_integration clean
	@echo "✅All checks are passed👌"

# -----------------------------TESTS------------------------
DC         := docker compose
DC_TEST    := ./test/docker-compose_test.yml
INT_TEST := ./test/go_test

test:
	go test -v -race $(INTERNAL)... $(PKG)...

test_integration:
	@$(DC) -f $(DC_TEST) down -v
	@$(DC) -f $(DC_TEST) up --build -d
	sleep 3
	@go test $(INT_TEST) -timeout 50s
	@$(DC) -f $(DC_TEST) down -v