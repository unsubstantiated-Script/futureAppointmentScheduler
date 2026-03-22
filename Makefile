SHELL := /bin/zsh

COMPOSE ?= docker compose
TEST_PKGS ?= ./internal/appointments

.PHONY: up down down-volumes test

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

down-volumes:
	$(COMPOSE) down -v

test:
	go test -v $(TEST_PKGS)
