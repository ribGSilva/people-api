SHELL := /bin/bash

# =====================================
# Variable

VERSION := 1.0
PROJECT_NAME := person-api

# =====================================
# Develop

debug:
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/cosmtrek/air@latest
	air -c .air.debug.toml

run:
	go install github.com/cosmtrek/air@latest
	air

tests:
	go test ./app/api/tests

tidy:
	go mod tidy
	go mod vendor

# =====================================
# Swagger

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init \
        --parseInternal \
        --parseDependency \
        --parseDepth 3 \
        --output app/api/docs \
        --dir app/api/

# =====================================
# Enviroment

env-build:
	-docker network create -d bridge $(PROJECT_NAME)
	-docker run -d -ti -p 27017:27017 --name mongo --network="$(PROJECT_NAME)" -e MONGO_INITDB_ROOT_USERNAME=person_app -e MONGO_INITDB_ROOT_PASSWORD=person_app_pass mongo
	-docker run -d -ti -p 6379:6379 --network="$(PROJECT_NAME)" --name redis redis

env-setup:
	-go run ./app/cmd/main.go schema create

env-up:
	-docker start mongo
	-docker start redis

env-down:
	-docker stop mongo
	-docker stop redis

env-clear:
	-docker rm mongo
	-docker rm redis
	-docker network rm $(PROJECT_NAME)

# =====================================
# Docker

docker-up: docker-build docker-run docker-logs

docker-reload: docker-stop docker-remove docker-build docker-run docker-logs

docker-down: docker-stop docker-remove docker-clear

docker-build:
	docker build -t $(PROJECT_NAME):$(VERSION) -f zarf/docker/Dockerfile .

docker-run:
	docker run -it -d -p 8080:8080 -p 4000:4000 --network="$(PROJECT_NAME)" \
 		-e MONGO_CONNECTION_URL="mongodb://person_app:person_app_pass@mongo:27017" \
 		-e REDIS_ADDRESS="redis:6379" \
 		--name $(PROJECT_NAME) $(PROJECT_NAME):$(VERSION)

docker-stop:
	docker stop $(PROJECT_NAME)

docker-stats:
	docker stats $(PROJECT_NAME)

docker-logs:
	docker logs -f --tail 20 $(PROJECT_NAME)

docker-remove:
	docker rm $(PROJECT_NAME)

docker-clear:
	docker rmi $(PROJECT_NAME):$(VERSION)
