
APP_NAME := note-app
DOCKER_REGISTRY ?= omegazyadav
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(APP_NAME)
DOCKER_TAG ?= latest
PORT := 8080

# Build the Go application
build:
	@echo "Building the Go application..."
	go build -o main .

# Run the application locally
run: build
	@echo "Running the application locally..."
	./main

# Build the Docker image
docker-build:
	@echo "Building the Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Push the Docker image to Docker Hub
docker-push: docker-build
	@echo "Pushing $(DOCKER_IMAGE):$(DOCKER_TAG) to Docker Hub..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run the Docker container
docker-run: docker-build
	@echo "Running the Docker container..."
	docker run --rm -p $(PORT):8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

deploy:
	kubectl apply -f k8s/

down: 
	kubectl delete -f k8s/

# Clean up build artifacts
clean:
	@echo "Cleaning up build artifacts..."
	rm -f main

# Help command
help:
	@echo "Available commands:"
	@echo "  build         Build the Go application"
	@echo "  run           Run the application locally"
	@echo "  docker-build  Build the Docker image"
	@echo "  docker-push   Build and push image to Docker Hub"
	@echo "  docker-run    Run the Docker container"
	@echo "  deploy        Build, push, and run the Docker container"
	@echo "  down          Stop the  container" 
	@echo "  clean         Clean up build artifacts"

.PHONY: build run docker-build docker-push docker-run clean help

