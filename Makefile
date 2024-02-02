include .env
export

.PHONY: run
run:
	@echo "Running the application"
	@go run . --limit 1

.PHONY: up
up:
	@echo "Starting the application"
	@docker-compose up -d

.PHONY: down
down:
	@echo "Stopping the application"
	@docker-compose down

.PHONY: up_build
up_build:
	@echo "Starting the application"
	@docker-compose up -d --build