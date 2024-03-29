#!/usr/bin/make -ef

NAME = go-space
MODULE_NAME = github.com/bahner/go-space
VAULT_TOKEN ?= space
NODE = $(NAME)-node

GO ?= go
PREFIX ?= /usr/local

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

default: clean tidy $(NAME)

init: go.mod tidy

go.mod:
	$(GO) mod init $(MODULE_NAME)

tidy: go.mod
	$(GO) mod tidy

node: $(NODE)

$(NODE): tidy
	$(GO) build -o $(NODE) ./cmd/node

clean:
	rm -f $(NAME) $(NODE)

console:
	docker-compose up -d
	docker attach go-space-pubsub_space_1

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

down:
	docker-compose down

image:
	docker build \
		-t $(IMAGE) \
		--build-arg "BUILD_IMAGE=$(BUILD_IMAGE)" \
		.

install: default
	sudo install -Dm755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)

run: clean $(NAME)
	./$(NAME)

up:
	docker-compose up -d --remove-orphans

vault:
	# docker-compose up -d vault
	vault server --dev -dev-root-token-id=$(VAULT_TOKEN) &

.PHONY: default init tidy build client serve install clean distclean
