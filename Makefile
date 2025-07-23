include .env

run-server:
	PORT=$(PORT) \
	POSTGRES_URI=$(POSTGRES_URI) \
	REDIS_URI=$(REDIS_URI) \
	AUTH_SECRET_KEY=$(AUTH_SECRET_KEY) \
	REPORTS_DIRECTORY=$(REPORTS_DIRECTORY) \
	PROMETHEUS_URL=$(PROMETHEUS_URL) \
	go run ./cmd/main.go

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-cli ./cmd/cli/main.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-server ./cmd/server/main.go

build-cli:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-cli ./cmd/cli/main.go

build-server:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-server ./cmd/server/main.go

