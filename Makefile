dependency:
	@echo ">> Downloading Dependencies"
	@go mod download

swag-init:
	@echo ">> Running swagger init"
	@swag init

run-api: dependency swag-init
	@echo ">> Running API Server"
	@go run main.go serve-http

migrate-up:
	@echo ">> Running Migrate Up"
	@migrate -path db/migrations -database "postgres://postgres:1235813@localhost:5433/users?sslmode=disable" up

migrate-down:
	@echo ">> Running Migrate down"
	@migrate -path db/migrations -database "postgres://postgres:1235813@localhost:5433/users?sslmode=disable" down