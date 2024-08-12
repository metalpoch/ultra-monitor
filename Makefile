dev-updater:
	CONFIG_JSON=./config.json go run ./update/cmd/main.go

build-updater:
	go build -o ./dist/olt-updater ./update/cmd/main.go

container-run:
	docker-compose up
