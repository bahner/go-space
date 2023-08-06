#!/usr/bin/make -ef

NAME = go-myspace-libp2p

init: go.mod go.sum

go.mod:
	go mod init $(NAME)

go.sum:
	go mod tidy
