dev-auth:
	CONFIG_JSON=./config.develop.json go run ./auth/cmd/main.go

dev-update:
	CONFIG_JSON=./config.develop.json go run ./update/cmd/main.go

build-updater:
	go build -o ./dist/olt-updater ./update/cmd/main.go

container-run:
	docker-compose up
