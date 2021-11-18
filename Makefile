PROJECT_NAME=cr-product
BUILD_VERSION=1.0.0

DOCKER_IMAGE=$(PROJECT_NAME):$(BUILD_VERSION)
GO_BUILD_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on

all: tidy vet build run

tidy:
	go mod tidy
vet:
	go vet ...
run:
	go run cmd/server/main.go
build:
	$(GO_BUILD_ENV) go build -v -o $(PROJECT_NAME)-$(BUILD_VERSION).bin cmd/server/main.go

compose_prod: docker
	cd deploy && BUILD_VERSION=$(BUILD_VERSION) docker-compose up  --build --force-recreate -d

docker_prebuild: build
	mkdir -p deploy/conf
	mv $(PROJECT_NAME)-$(BUILD_VERSION).bin deploy/$(PROJECT_NAME).bin; \
	cp -R conf deploy/;
docker_build:
	cd deploy; \
	docker build -t $(DOCKER_IMAGE) .;

docker_postbuild:
	cd deploy; \
	rm -rf $(PROJECT_NAME).bin 2> /dev/null;\
	rm -rf conf 2> /dev/null;
docker: docker_prebuild docker_build docker_postbuild

mock:
	go generate -x -run="mockgen" ./...
