BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER := postgres-calendar

#Rabbit
RABBIT_USER ?= rabbit
RABBIT_PASSWORD ?= password
RABBIT_PORT ?= 5672
RABBIT_PORT_UI ?= 15672
RABBIT_CONTAINER := rabbitmq

up:
	docker compose -f deployments/docker-compose.yaml up -d --build

down:
	docker compose -f deployments/docker-compose.yaml down

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.yaml

run-api: build
	$(BIN) -config ./configs/config.yaml

run-scheduler: build
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml

run-sender: build
	$(BIN_SENDER) -config ./configs/sender_config.yaml

run-all: build
	trap 'kill 0' SIGINT SIGTERM EXIT; \
	$(MAKE) run-api & \
	$(MAKE) run-scheduler & \
	$(MAKE) run-sender & \
	wait

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

api-test:
	go test -v ./internal/server/http/...

lint:
	go mod download
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5 run ./...

migration:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir migrations create $(NAME) sql

generate:
	mkdir -p pkg/api
	protoc --proto_path=api \
		--proto_path=internal/proto \
		--go_out=pkg/api --go_opt=paths=source_relative \
		--go-grpc_out=pkg/api --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pkg/api --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=pkg/api --openapiv2_opt=logtostderr=true \
		api/EventService.proto

load-tests:
	k6 run load_tests/crud.js

.PHONY: build run run-api run-scheduler run-sender run-all build-img run-img version test api-test lint migration generate load-tests up down
