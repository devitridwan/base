APP_NAME=payslip-service

dep:
	@echo ">> Downloading Dependencies"
	@go mod download

run-api:
	@echo ">> Running API Server"
	@go run main.go serve-http

swag-init:
	@echo ">> Running swagger init"
	@swag init

remock:
	#https://github.com/vektra/mockery

	@echo ">> Mocking Files"
	@mockery 

run-test: dep
	@echo ">> Running Test"
	@go test -v -cover -count=1 -failfast -covermode=atomic ./...