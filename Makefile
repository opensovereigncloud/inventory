BIN_NAME_INVENTORY = "inventory"
BIN_NAME_BENCHMARK = "benchmark"

DOCKER_REGISTRY = "localhost:5000"
DOCKER_IMAGE_NAME = "inventory"
DOCKER_IMAGE_TAG = "latest"

GOPRIVATE ?= "github.com/onmetal/*"

.PHONY: all
all: compile

.PHONY: compile
compile: fmt vet
	go build -o dist/$(BIN_NAME_INVENTORY) cmd/inventory/main.go
	go build -o dist/$(BIN_NAME_BENCHMARK) cmd/benchmark/main.go
	cp -rf res/ dist/

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: dl-pciids
dl-pciids:
	curl https://pci-ids.ucw.cz/v2.2/pci.ids --output ./res/pci.ids

.PHONY: docker-build
docker-build:
	docker build . -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) --build-arg GOPRIVATE=$(GOPRIVATE) --build-arg GIT_USER=$(GIT_USER) --build-arg GIT_PASSWORD=$(GIT_PASSWORD)

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

.PHONY: clean
clean:
	rm -rf ./dist/
