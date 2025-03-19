dev-auth:
	CONFIG_JSON=./config.json go run ./auth/cmd/main.go

dev-core:
	CONFIG_JSON=./config.json go run ./core/cmd/main.go

dev-report:
	CONFIG_JSON=./config.json go run ./report/cmd/main.go


build:
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-auth ./auth/cmd/main.go && echo -e "\e[1;32mcreated\e[0m was created the binary olt-blueprint-auth"
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-report ./report/cmd/main.go && echo -e "\e[1;32mcreated\e[0m was created the binary olt-blueprint-report"
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-core ./core/cmd/main.go && echo -e "\e[1;32mcreated\e[0m was created the binary olt-blueprint-core"
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./dist/olt-blueprint-measurement ./measurement/cmd/main.go && echo -e "\e[1;32mcreated\e[0m was created the binary olt-blueprint-measurement"

measurement-build-cli:
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./measurement/dist/olt-blueprint measurement/cmd/main.go && echo -e "\e[1;32mcreated\e[0m binary was measurement/dist/olt-blueprint"

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
	cd smart
	docker build . -t smart  -f ./smart/dockerfile

container-smart-run:
	docker run --rm -p 3003:3003 --name olt-blueprint-smart olt-blueprint-smart
	
