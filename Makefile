# Chprobe Project Management Makefile

# Basic configuration
SystemName="chprobe"
GoVersion=$(shell go version | awk '{print $$3, $$4}')
BuildTime=$(shell date "+%Y-%m-%d %H:%M:%S")
GitCommitID=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
Version="0.0.1_base"

# Build mode
VER=debug

# Set Gin log mode based on build mode
ifeq ($(VER),debug)
	GinLogMode="debug"
else
	GinLogMode="release"
endif

# Build parameters
LDFLAGS="\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.GoVersion=$(GoVersion)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.SystemName=$(SystemName)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.CommitID=$(GitCommitID)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.BuildTime=$(BuildTime)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.Version=$(Version)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.GinLogMode=${GinLogMode}'\
"

# Define start-build function
define start-build=
	@echo $(1)
	@make $(1)
endef

# Target definitions
.PHONY: help stop all client server web prepare revert docker-build docker-run docker-stop release release-client release-server

# Display help information
help:
	@echo "Chprobe Project Management Tool"
	@echo "Usage: make [target]"
	@echo ""
	@echo "Main Targets:"
	@echo "  help          - Display this help information"
	@echo "  build         - Build all components (release mode)"
	@echo "  build-client  - Build client (release mode)"
	@echo "  build-server  - Build server (release mode)"
	@echo "  clean         - Clean all build artifacts"
	@echo "  stop          - Clean running processes and build artifacts"
	@echo "  prepare       - Prepare database environment"
	@echo "  revert        - Revert prepare actions (stop and remove containers)"
	@echo "  all           - Clean all processes and run all components"
	@echo ""
	@echo "Docker Targets:"
	@echo "  docker-build  - Build server Docker image (uses pre-built binaries)"
	@echo "  docker-run    - Run server Docker container"
	@echo "  docker-stop   - Stop and remove server Docker container"
	@echo ""
	@echo "Component Targets:"
	@echo "  client        - Clean client process, build and run client"
	@echo "  server        - Clean server process, build and run server"
	@echo "  web           - Clean web process, build and run web frontend"
	@echo ""
	@echo "Build Modes:"
	@echo "  VER=debug     - Debug mode (default)"
	@echo "  VER=release   - Release mode"
	@echo ""
	@echo "Examples:"
	@echo "  make all              # Run all components (debug mode)"
	@echo "  make server           # Run only server (debug mode)"
	@echo "  make all VER=release  # Run all components (release mode)"
	@echo "  make stop             # Clean all running processes"
	@echo "  make prepare          # Prepare database environment"
	@echo "  make revert           # Revert prepare actions"
	@echo "  make build            # Build all components (release mode)"
	@echo "  make build-client     # Build client (release mode)"
	@echo "  make build-server     # Build server (release mode)"
	@echo "  make docker-build     # Build server Docker image (uses pre-built binaries)"
	@echo "  make docker-run       # Run server Docker container"
	@echo "  make docker-stop      # Stop and remove server Docker container"



# Clean running processes
stop:
	@echo "Cleaning running processes..."
	@# Stop client process
	@ps aux | grep chprobe_client | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Stop server process
	@ps aux | grep chprobe_server | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Stop web frontend process (npm run dev)
	@ps aux | grep "npm run dev" | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Stop all chprobe-related processes
	@ps aux | grep chprobe | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@echo "Process cleanup completed!"

	@echo "Cleaning all build artifacts..."
	@rm -rf ./chprobe_client/bin 2>/dev/null || true
	@rm -rf ./chprobe_server/bin 2>/dev/null || true
	@rm -rf ./chprobe_web/dist 2>/dev/null || true
	@echo "Cleanup completed!"

# Prepare database environment
prepare:
	@echo "Starting to prepare database environment..."
	@# Execute server's prepare script
	@cd ./chprobe_server && bash ./scripts/prepare.sh
	@echo "Database environment preparation completed!"

# Revert prepare actions (stop and remove containers)
revert:
	@echo "Starting to revert prepare actions..."
	@echo "Stopping and removing mysql container..."
	@sudo docker stop mysql 2>/dev/null || true
	@sudo docker rm mysql 2>/dev/null || true
	@echo "Stopping and removing redis container..."
	@sudo docker stop redis 2>/dev/null || true
	@sudo docker rm redis 2>/dev/null || true
	@echo "Containers stopped and removed, environment reverted!"

# Run all components
all:
	@echo "Starting to run all components..."
	@echo "Note: This command will run all components in the background"
	@# Clean all processes
	@make stop
	@# Build and run all components
	@make server
	@sleep 2
	@make web
	@sleep 2
	@make client
	@echo "All components started!"

# Client related target
client:
	@echo "Starting to build and run client..."
	@# Clean client process
	@ps aux | grep chprobe_client | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Build client
	@cd ./chprobe_client && mkdir -p bin && go build -o bin/chprobe_client ./cmd/core/main.go
	@echo "Client build completed!"
	@# Run client
	@cd ./chprobe_client && ./bin/chprobe_client &
	@echo "Client started!"

# Server related target
server:
	@echo "Starting to build and run server..."
	@# Clean server process
	@ps aux | grep chprobe_server | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Build server
	@cd ./chprobe_server && mkdir -p bin && go build -ldflags $(LDFLAGS) -o bin/chprobe_server ./cmd/main.go
	@echo "Server build completed!"
	@# Run server
	@cd ./chprobe_server && ./bin/chprobe_server &
	@echo "Server started!"

# Web frontend related target
web:
	@echo "Starting to build and run web frontend..."
	@# Clean web frontend process
	@ps aux | grep "npm run dev" | grep -v grep | awk '{print $$2}' | xargs kill -9 2>/dev/null || true
	@# Build web frontend
	@cd ./chprobe_web && npm install && npm run build
	@echo "Web frontend build completed!"
	@# Run web frontend
	@cd ./chprobe_web && npm run dev &
	@echo "Web frontend started!"

# Build targets

# Build all components in release mode
release:
	@echo "Building all components in release mode..."
	@make release-client
	@make release-server
	@echo "All components built successfully!"

# Build client in release mode
release-client:
	@echo "Building client in release mode..."
	@cd ./chprobe_client && mkdir -p bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o bin/chprobe_client ./cmd/core/main.go
	@echo "Client built successfully!"

# Build server in release mode
release-server:
	@echo "Building server in release mode..."
	@cd ./chprobe_server && mkdir -p bin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o bin/chprobe_server ./cmd/main.go
	@echo "Server built successfully!"

# Docker related targets

# Build server Docker image (uses pre-built binaries)
docker-build:
	@echo "Starting to build server Docker image..."
	@# Ensure binaries are built
	@make release
	@sudo docker build -t chprobe-server -f Dockerfile .
	@echo "Server Docker image built successfully!"

# Run server Docker container
docker-run:
	@make docker-build
	@echo "Starting to run server Docker container..."
	@# Stop and remove existing container if it exists
	@sudo docker stop chprobe-server 2>/dev/null || true
	@sudo docker rm chprobe-server 2>/dev/null || true
	@# Run new container
	@sudo docker run --name chprobe-server --restart always -p 32000:32000 -p 32001:32001 -d chprobe-server
	@echo "Server Docker container started successfully!"

# Stop and remove server Docker container
docker-stop:
	@echo "Stopping and removing server Docker container..."
	@sudo docker stop chprobe-server 2>/dev/null || true
	@sudo docker rm chprobe-server 2>/dev/null || true
	@echo "Server Docker container stopped and removed!"

