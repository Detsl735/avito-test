APP_NAME=pr-service
APP_PATH=./cmd/app

DOCKER_COMPOSE=docker compose -f deployments/docker-compose.yml


build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) $(APP_PATH)

run:
	@echo "Running $(APP_NAME) locally..."
	go run $(APP_PATH)

test:
	@echo "Running tests..."
	go test ./... -v -count=1

lint:
	@echo "Running golangci-lint..."
	golangci-lint run

tidy:
	@echo "Tidying modules..."
	go mod tidy


docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

up:
	@echo "Starting services via docker-compose..."
	$(DOCKER_COMPOSE) up --build

up-d:
	@echo "Starting services in background..."
	$(DOCKER_COMPOSE) up -d --build

down:
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f app


clean:
	@echo "Cleaning..."
	rm -rf bin

.PHONY: build run test lint tidy docker-build up up-d down logs clean
