#!/usr/bin/make -ef

NAME = myspace-pubsub

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

$(NAME): tidy
	$(GO) build -o $(NAME)

clean:
	rm -f $(NAME)

distclean: clean
	rm -f $(shell git ls-files --exclude-standard --others)

image:
	docker build \
		-t $(IMAGE) \
		--build-arg "BUILD_IMAGE=$(BUILD_IMAGE)" \
		.

run: clean $(NAME)
	./$(NAME)

.PHONY: default init tidy build client serve install clean distclean

install: $(NAME)
	install -Dm755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)

