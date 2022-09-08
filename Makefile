APP_NAME=golang-grpc-base-project
APP_VERSION=latest
DOCKER_REGISTRY=registry.gitlab.com/xdorro/registry

docker.build:
	docker build -t $(DOCKER_REGISTRY)/$(APP_NAME):$(APP_VERSION) .

docker.push:
	docker push $(DOCKER_REGISTRY)/$(APP_NAME):$(APP_VERSION)

docker.dev: docker.build docker.push

wire.gen:
	wire ./...

lint.run:
	golangci-lint run --fast ./...

go.install:
	go install github.com/google/wire/cmd/wire@latest

	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

go.gen: wire.gen

go.lint: lint.run

go.get:
	go get -u ./...

go.tidy:
	go mod tidy

go.test:
	go test ./...

