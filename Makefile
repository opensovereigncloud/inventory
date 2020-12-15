BIN_NAME = "inventory"

DOCKER_REGISTRY = "localhost:5000"
DOCKER_IMAGE_NAME = "inventory"
DOCKER_IMAGE_TAG = "latest"

.PHONY: all
all: compile

.PHONY: compile
compile: fmt vet
	go build -o bin/$(BIN_NAME) cmd/inventory/main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: docker-build
docker-build:
	docker build . -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: clean
clean:
	rm -rf ./bin/