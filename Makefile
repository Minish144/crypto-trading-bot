#!make
include .env
export $(shell sed 's/=.*//' .env)

.PHONY:run
run:
	go run ./main.go

.PHONY:run-race
run-race:
	go run --race ./main.go

.PHONY:generate
generate:
	go generate ./...

.PHONY:lint-latest
lint-latest:
	@docker run --rm --pull never -v $(PWD):/app -w /app golangci/golangci-lint:v1.44-alpine golangci-lint run -v --timeout 5m

.PHONY:lint
lint:
	golangci-lint run

.PHONY:test
test:
	go test -v ./...

test-richgo:
	@RICHGO_FORCE_COLOR=1 richgo test ./...

.PHONY:gosec
gosec:
	@docker run --rm -v $(PWD):/app -w /app ${DOCKER_PROXY}/securego/gosec /app/...

path=$(dir ${file})
filename=$(notdir ${file})
mocks: # 'make mocks file=<filepath with interface>'
	@mockgen -source=${file} -destination=${path}mocks/${filename} -package mocks
	@echo "Generated mock '${path}mocks/${filename}'"

file=$(file)
gotests: # 'make gotests file=<filepath>'
	@gotests -w -exported ${file}

.DEFAULT_GOAL=run
