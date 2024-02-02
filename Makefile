include .env
export

.PHONY: run
run:
	@echo "Running the application"
	@go run . --limit 1