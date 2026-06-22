dependency:
	@echo ">> Downloading Dependencies"
	@go mod download

swag-init:
	@echo ">> Running swagger init"
	@swag init

run-api: dependency swag-init
	@echo ">> Running API Server"
	@go run main.go serve-http

run-grpc: dependency
	@echo ">> Running gRPC Server"
	@go run main.go serve-grpc


migrate-up:
	@echo ">> Running Migrate Up"
	@migrate -path db/migrations -database "postgres://hanifabyana:@localhost:5432/garnet?sslmode=disable" up

migrate-down:
	@echo ">> Running Migrate down"
	@migrate -path db/migrations -database "postgres://hanifabyana:@localhost:5432/garnet?sslmode=disable" down

proto:
	@echo ">> Generating proto, gRPC stubs, grpc-gateway, and OpenAPI v2 swagger"
	@buf generate

proto-gen:
	@echo ">> Generating gRPC proto files"
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/handler/grpc/proto/user.proto


remock:
	#https://github.com/vektra/mockery
	@echo ">> Mock Repositories"
	@mockery --all --dir ./internal/domain/repository --output ./internal/domain/repository/mocks_repository --outpkg mocks_repository

	@echo ">> Mock UseCases"
	@mockery --all --dir ./internal/usecases --output ./internal/usecases/mocks_usecases --outpkg mocks_usecases

	@echo ">> Mock Interfaces"
	@mockery --all --recursive --dir ./internal/interfaces --output ./internal/interfaces/mocks_interfaces --outpkg mocks_interfaces

	@echo ">> Mock Infra"
	@mockery --all --recursive --dir ./internal/infrastructures --output ./internal/infrastructures/mocks_infrastructures --outpkg mocks_infrastructures