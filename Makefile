APP_NAME := autotransport
BINARY_DIR := bin
DOCKER_IMAGE := autotransport:latest
DOCKER_COMPOSE_FILE := docker-compose.yml

clean:
	@echo "Cleaning up binaries..."
	rm -rf ./$(BINARY_DIR)

build_app: clean
	@echo "Building Go application..."
	mkdir -p $(BINARY_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./$(BINARY_DIR)/$(APP_NAME) ./cmd/main.go

build_image: build_app
	@echo "Building Docker image with PostgreSQL Alpine..."
	docker build -t $(DOCKER_IMAGE) -f ./Dockerfile .

export_image: build_image
	@echo "Exporting Docker image to a tarball..."
	docker save -o ./$(DOCKER_IMAGE).tar $(DOCKER_IMAGE)

clean_containers:
	@echo "Stopping and removing existing containers..."
	docker compose -f ./$(DOCKER_COMPOSE_FILE) down

build_containers: clean_containers build_image
	@echo "Building containers using Docker Compose..."
	docker compose -f ./$(DOCKER_COMPOSE_FILE) up --build -d

build_containers_local: clean_containers build_image
	@echo "Building containers using Docker Compose..."
	docker compose -f ./$(DOCKER_COMPOSE_FILE) up --build

run_local: build_app
	@echo "Running application locally..."
	./$(BINARY_DIR)/$(APP_NAME)

all: build_app  build_containers
all_local: build_app  build_containers_local

.PHONY: clean build_app build_image build_containers clean_containers export_image run_local all