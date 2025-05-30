PROJECT_NAME := skypiea
WORKDIR = /app
CHART_PATH := $(CURRENT_DIR)/helm/skypiea
DEFAULT_CONFIG = config.example.yml
CONFIG_FLAG ?= --config=$(DEFAULT_CONFIG)

CURRENT_DIR ?= $(shell pwd)
TIMESTAMP := $(shell date +%s)

run:
	go run ./cmd/server/main.go $(CONFIG_FLAG)

run-worker:
	go run ./cmd/worker/main.go $(CONFIG_FLAG)

build-server:
	go build -o bin/server ./cmd/server

build-worker:
	go build -o bin/worker ./cmd/worker

.PHONY: build
build: build-server build-worker

migratedb-up:
	go run ./cmd/server/main.go migration up $(CONFIG_FLAG)

migratedb-down:
	go run ./cmd/server/main.go migration down $(CONFIG_FLAG)

lint:
	golangci-lint run -v

docker-build-server:
	docker build -f ./docker/server.Dockerfile -t skypiea-ai-server:latest .

docker-run-server:
	docker run -d --rm --name skypiea-ai-server -p 8080:8080 -e DATABASE_POSTGRESQL_HOST=host.docker.internal skypiea-ai-server:latest

docker-build-worker:
	docker build -f ./docker/worker.Dockerfile -t skypiea-ai-worker:latest .

docker-run-worker: #no need to use
	@echo 'Houston, we are launching'
	@# docker run -d --rm --name skypiea-ai-worker skypiea-ai-worker:latest

docker-compose-build-up:
	docker-compose up --build

docker-compose-db-only:
	docker-compose up postgresdb -d

TEST_DIRS = ./internal/... ./cmd/... ./pkg/...

.PHONY: test coverage

test:
	@go test $(TEST_DIRS)

# Run tests with coverage and print a summary
coverage:
	@go test $(TEST_DIRS) -coverprofile=coverage.out
	@go tool cover -func=coverage.out