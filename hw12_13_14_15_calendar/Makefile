BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DATABASE_DSN := "postgresql://postgres:postgres@localhost:54321/calendar?sslmode=disable"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)_scheduler -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN)_sender -ldflags "$(LDFLAGS)" ./cmd/sender

run: build
	$(BIN)

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

integration-tests:
	set -e ;\
	docker-compose up --build -d ;\
	test_status_code=0 ;\
	docker-compose run tests go test github.com/koind/avito_otus_hw/hw12_13_14_15_calendar/tests/integrations || test_status_code=$$? ;\
	docker-compose down ;\
	exit $$test_status_code ;

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	goose --dir=migrations postgres $(DATABASE_DSN) up

generate:
	go generate ./...

up:
	docker-compose up -d

down:
	docker-compose down

.PHONY: build run build-img run-img version test lint