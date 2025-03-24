dev-auth:
	CONFIG_JSON=./config.json go run ./auth/cmd/main.go

dev-core:
	CONFIG_JSON=./config.json go run ./core/cmd/main.go

dev-report:
	CONFIG_JSON=./config.json go run ./report/cmd/main.go

build:
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-auth ./auth/cmd/main.go
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-report ./report/cmd/main.go
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-core ./core/cmd/main.go
	CGO_ENABLED=0 go build -o ./dist/olt-blueprint-measurement ./measurement/cmd

build-img-containers:
	docker build . -t olt-blueprint-auth --progress=plain -f ./auth/Dockerfile
	docker build . -t olt-blueprint-core --progress=plain -f ./core/Dockerfile
	docker build . -t olt-blueprint-front --progress=plain -f ./client/Dockerfile
	docker build . -t olt-blueprint-report --progress=plain -f ./report/Dockerfile
	docker build . -t olt-blueprint-measurement --progress=plain -f ./measurement/Dockerfile

container-run:
	docker-compose up

container-auth: 
	docker build . -t olt-blueprint-auth --progress=plain -f ./auth/Dockerfile

container-core: 
	docker build . -t olt-blueprint-core --progress=plain -f ./core/Dockerfile

container-report: 
	docker build . -t olt-blueprint-report --progress=plain -f ./report/Dockerfile

container-measurement:
	docker build . -t olt-blueprint-cli --progress=plain -f ./measurement/Dockerfile

container-measurement-cli:
	docker run --rm -v ./config.json:/app/config.json --name olt-blueprint-cli -e CONFIG_JSON='/app/config.json' olt-blueprint-cli

container-smart:
	docker build . -t olt-blueprint-smart --progress=plain -f ./smart/Dockerfile

container-smart-run:
	docker run --rm -p 3003:3003 --network host --name olt-blueprint-smart olt-blueprint-smart

