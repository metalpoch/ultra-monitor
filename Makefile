dev-auth:
	CONFIG_JSON=./config.json go run ./auth/cmd/main.go

dev-traffic:
	CONFIG_JSON=./config.json go run ./traffic/cmd/main.go

measurement-build-cli:
	go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o ./measurement/dist/olt-blueprint measurement/cmd/main.go && echo -e "\e[1;32mcreated\e[0m binary was measurement/dist/olt-blueprint"

container-run:
	docker-compose up

container-measurement:
	docker build . -t olt-blueprint-cli --progress=plain -f ./measurement/Dockerfile

container-measurement-cli:
	docker run --rm -v ./config.json:/app/config.json --name olt-blueprint-cli -e CONFIG_JSON='/app/config.json' olt-blueprint-cli

container-chatbot:
	cd chatbot
	docker build . -t chatbot  -f ./chatbot/dockerfile
	docker run -p 8000:8000 --name chatbot-olt --network olt-blueprint_chatbot-network chatbot  
	