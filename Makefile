SHELL := /bin/bash

ROOT_DIR := $(CURDIR)
BACKEND_DIR := $(ROOT_DIR)/backend
FRONTEND_DIR := $(ROOT_DIR)/web

COMPOSE ?= docker compose
BUN ?= bun
GO ?= go
AIR ?= air

POSTGRES_DSN ?= postgres://postgres:postgres@localhost:5432/your_next_game?sslmode=disable

.PHONY: help up-containers down-containers logs-containers deps setup migrate backend-dev frontend-dev dev

help:
	@printf '%s\n' \
		'Available targets:' \
		'  make up-containers   - start postgres + otel collector' \
		'  make backend-dev     - run backend in dev mode with air' \
		'  make frontend-dev    - run Next.js in dev mode' \
		'  make dev             - run backend and frontend together' \
		'  make migrate         - apply backend migrations' \
		'  make setup           - install deps, start containers, and run migrations' \
		'  make down-containers - stop compose services'

up-containers:
	$(COMPOSE) up -d postgres otel-collector loki grafana

down-containers:
	$(COMPOSE) down

logs-containers:
	$(COMPOSE) logs -f postgres otel-collector loki grafana

deps:
	cd "$(FRONTEND_DIR)" && $(BUN) install
	cd "$(BACKEND_DIR)" && $(GO) mod download

setup:
	$(MAKE) deps
	$(MAKE) migrate

migrate: up-containers
	cd "$(BACKEND_DIR)" && POSTGRES_DSN="$(POSTGRES_DSN)" $(GO) run ./cmd/migrate

backend-dev: up-containers
	cd "$(BACKEND_DIR)" && POSTGRES_DSN="$(POSTGRES_DSN)" go run ./cmd/api/main.go

frontend-dev:
	cd "$(FRONTEND_DIR)" && $(BUN) dev

dev: up-containers
	@trap 'kill 0' INT TERM EXIT; \
		cd "$(BACKEND_DIR)" && POSTGRES_DSN="$(POSTGRES_DSN)" go run ./cmd/api/main.go & \
		cd "$(FRONTEND_DIR)" && $(BUN) dev & \
		wait
