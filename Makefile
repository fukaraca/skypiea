PROJECT_NAME := skypiea
WORKDIR = /app
CHART_PATH := $(CURRENT_DIR)/helm/skypiea
DEFAULT_CONFIG = config.example.yml
CONFIG_FLAG ?= --config=$(DEFAULT_CONFIG)

CURRENT_DIR ?= $(shell pwd)
TIMESTAMP := $(shell date +%s)

run-server:
	go run ./cmd/server/main.go $(CONFIG_FLAG)

run-worker:
	go run ./cmd/worker/main.go $(CONFIG_FLAG)

build-server:
	go build -o bin/server ./cmd/server

build-worker:
	go build -o bin/worker ./cmd/worker

.PHONY: build
build: build-server build-worker