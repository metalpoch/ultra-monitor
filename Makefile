include .env

build:
	make build-server
	make build-web

build-server:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/gestor-ultra ./cmd/main.go

build-web:
	cd web/ && npm run build

run-web:
	cd web/ && npm run dev

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

start-server: build
	PORT=$(PORT) \
	POSTGRES_URI=$(POSTGRES_URI) \
	REDIS_URI=$(REDIS_URI) \
	AUTH_SECRET_KEY=$(AUTH_SECRET_KEY) \
	REPORTS_DIRECTORY=$(REPORTS_DIRECTORY) \
	PROMETHEUS_URL=$(PROMETHEUS_URL) \
	./dist/gestor-ultra server
