build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-cli ./cmd/cli/main.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-server ./cmd/server/main.go

build-cli:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-cli ./cmd/cli/main.go

build-server:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ./dist/ultra-server ./cmd/server/main.go

