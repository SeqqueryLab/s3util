#!make
include .env
export

test:
	env

hello:
	@echo '${NAME} ${VERSION} ${ENVIRONMENT}'

run:
	@go run ./cmd

build:
	@echo 'Building ${NAME} ${VERSION}' 
	@export PATH=$(go env GOPATH)/bin:$PATH;
	@swag init -g cmd/server/main.go > /dev/null
	@go build ./cmd

docker:
	@echo 'Building Docker image of ${NAME} ${VERSION}'
	@docker build -t ${NAME}:${VERSION} .

all: hello run