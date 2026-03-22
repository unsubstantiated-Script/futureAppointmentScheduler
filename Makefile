SHELL := /bin/zsh

COMPOSE ?= docker compose

.PHONY: up down down-volumes

up:
	$(COMPOSE) up --build

down:
	$(COMPOSE) down

down-volumes:
	$(COMPOSE) down -v

