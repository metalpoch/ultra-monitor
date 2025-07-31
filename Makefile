include .env

run-server:
	PORT=$(PORT) \
	POSTGRES_URI=$(POSTGRES_URI) \
	REDIS_URI=$(REDIS_URI) \
	AUTH_SECRET_KEY=$(AUTH_SECRET_KEY) \
	REPORTS_DIRECTORY=$(REPORTS_DIRECTORY) \
	PROMETHEUS_URL=$(PROMETHEUS_URL) \
	go run ./cmd/main.go server

run-scan:
	PORT=$(PORT) \
	POSTGRES_URI=$(POSTGRES_URI) \
	REDIS_URI=$(REDIS_URI) \
	AUTH_SECRET_KEY=$(AUTH_SECRET_KEY) \
	REPORTS_DIRECTORY=$(REPORTS_DIRECTORY) \
	PROMETHEUS_URL=$(PROMETHEUS_URL) \
	go run ./cmd/main.go scan

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/gestor-ultra ./cmd/cli/main.go
