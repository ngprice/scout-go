APP_NAME := scout-go
BUILD_DIR := .
BUILD_PATH := $(BUILD_DIR)/main.go

.PHONY: all deps build run test

all: build

deps:
	@echo "Getting dependencies..."
	go get ./...
	@echo "Dependencies downloaded."

build: deps
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) $(BUILD_PATH)
	@echo "Build complete. Executable: ./${APP_NAME}"

test:
	@echo "Running tests..."
	go test -v ./...

run: build
	@echo "Running $(APP_NAME)..."
	./$(APP_NAME)
