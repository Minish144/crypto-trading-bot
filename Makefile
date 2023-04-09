#!make
include .env
export $(shell sed 's/=.*//' .env)

.PHONY:run
run:
	go run ./main.go

