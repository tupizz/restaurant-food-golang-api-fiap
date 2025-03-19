include .env

# Makefile

SHELL := /bin/bash

APP_NAME := fastfood-golang
GO := go
GOPATH := $(shell go env GOPATH)
SWAG := $(GOPATH)/bin/swag
AIR := $(GOPATH)/bin/air
MIGRATE := $(GOPATH)/bin/migrate
BINARY_NAME := $(APP_NAME)

export PATH := $(GOPATH)/bin:$(PATH)

.PHONY: help all build run run-air migration migrate-up migrate-down swag-init docker-up docker-down test clean install-tools setup sqlc

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  help           Show this help message"
	@echo "  build          Build the application binary"
	@echo "  run            Run the application"
	@echo "  run-air        Run the application with live reloading using Air"
	@echo "  migration      Create a new database migration"
	@echo "  migrate-up     Apply all database migrations"
	@echo "  migrate-down   Revert the last database migration"
	@echo "  swag-init      Generate Swagger documentation"
	@echo "  docker-up      Start Docker containers with docker-compose up --build"
	@echo "  docker-down    Stop Docker containers"
	@echo "  test           Run tests"
	@echo "  clean          Remove built binaries"
	@echo "  install-tools  Install required tools (swag, air, migrate)"
	@echo "  setup          Install tools, generate docs, and run migrations"
	@echo "  sqlc           Generate SQLC"

all: build

build:
	@echo "Building application..."
	$(GO) build -o bin/$(APP_NAME) ./cmd/main.go

run:
	@echo "Running application..."
	$(GO) run ./cmd/main.go

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o bin/$(BINARY_NAME)-debug ./cmd/main.go

run-air:
	@echo "Running application with Air..."
	$(AIR)

migration:
	@echo "Creating migration files for '$(filter-out $@,$(MAKECMDGOALS))'..."
	@migrate create -ext sql -dir ./database/migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@echo "Applying database migrations..."
	$(MIGRATE) -database $(DATABASE_URL) -path ./database/migrations up

migrate-down:
	@echo "Reverting database migrations..."
	$(MIGRATE) -database $(DATABASE_URL) -path ./database/migrations down 1

swag-init:
	@echo "Generating Swagger documentation..."
	$(SWAG) init -g cmd/main.go -o ./swagger

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up --build

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

test:
	@echo "Running tests..."
	$(GO) test -v ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/*

sqlc:
	@echo "Generating SQLC..."
	sqlc generate

install-tools:
	@echo "Installing tools..."
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/go-delve/delve/cmd/dlv@latest
	$(GO) install github.com/air-verse/air@latest
	$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

setup: install-tools swag-init migrate-up
