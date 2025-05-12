.PHONY: run-app db-up db-down integration-up integration-down integration-test run-test

include .env
export

run-app:
	@echo "Запуск приложения локально"
	go build -o bin/main_service ./cmd/app/main.go
	
	CONFIG_PATH=./config/config.local.yaml ./bin/main_service

# Запуск PostgreSQL в Docker с параметрами из .env
db-up:
	@echo "Запуск контейнера PostgreSQL..."
	docker run --rm --name local-postgres \
	  -e POSTGRES_USER=${POSTGRES_USER} \
	  -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
	  -e POSTGRES_DB=${POSTGRES_DB} \
	  -p ${POSTGRES_PORT}:5432 \
	  -d postgres:17

# Остановка контейнера PostgreSQL
db-down:
	@echo "Остановка контейнера PostgreSQL..."
	docker stop local-postgres

integration-up:
	docker compose -f docker-compose-integration-test.yaml up -d

integration-test: integration-up
	go test -v ./integration-test/...
	sleep 2
	make integration-down

integration-down:
	docker compose -f docker-compose-integration-test.yaml down

run-test: integration-up
	go test -v ./...
	sleep 1
	make integration-down