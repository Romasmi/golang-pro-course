BIN := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG_CALENDAR := calendar:latest
DOCKER_IMG_SCHEDULER := scheduler:latest
DOCKER_IMG_SENDER := sender:latest
K8S_NAMESPACE := calendar-app

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
		-t $(DOCKER_IMG_CALENDAR) \
		-f deployments/calendar.Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f deployments/scheduler.Dockerfile .
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f deployments/sender.Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG_CALENDAR)

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

integration-tests:
	docker compose -f deployments/docker-compose.integration.yaml up --build --exit-code-from integration-tests; \
	RET=$$?; \
	docker compose -f deployments/docker-compose.integration.yaml down; \
	exit $$RET

k8s-install:
	helm upgrade --install calendar deployments/calendar-chart --namespace $(K8S_NAMESPACE) --create-namespace

k8s-uninstall:
	helm uninstall calendar --namespace $(K8S_NAMESPACE)

k8s-load-images: build-img
	minikube image load $(DOCKER_IMG_CALENDAR) $(DOCKER_IMG_SCHEDULER) $(DOCKER_IMG_SENDER)

k8s-forward-api:
	kubectl port-forward -n $(K8S_NAMESPACE) svc/calendar-calendar-chart-calendar 8888:80

k8s-forward-db:
	kubectl port-forward -n $(K8S_NAMESPACE) svc/calendar-calendar-chart-postgres $(POSTGRES_PORT):5432

k8s-status:
	kubectl get pods -n $(K8S_NAMESPACE)

k8s-logs-api:
	kubectl logs -f -n $(K8S_NAMESPACE) -l app=calendar

k8s-logs-scheduler:
	kubectl logs -f -n $(K8S_NAMESPACE) -l app=scheduler

k8s-logs-sender:
	kubectl logs -f -n $(K8S_NAMESPACE) -l app=sender

k8s-logs:
	kubectl logs -f -n $(K8S_NAMESPACE) -l 'app in (calendar, scheduler, sender)' --all-containers=true --prefix

.PHONY: build run run-api run-scheduler run-sender run-all build-img run-img version test api-test lint migration generate load-tests up down integration-tests k8s-install k8s-uninstall k8s-load-images k8s-forward-api k8s-forward-db k8s-status k8s-logs-api k8s-logs-scheduler k8s-logs-sender k8s-logs
