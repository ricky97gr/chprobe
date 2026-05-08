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
	@echo "  build         - 编译所有产物 (前端+后端+客户端)"
	@echo "  build-web     - 仅编译前端"
	@echo "  build-client  - 仅编译客户端Agent"
	@echo "  build-server  - 仅编译服务端"
	@echo "  clean         - Clean all build artifacts"
	@echo "  stop          - Clean running processes and build artifacts"
	@echo "  prepare       - Prepare database environment"
	@echo "  revert        - Revert prepare actions (stop and remove containers)"
	@echo "  all           - Clean all processes and run all components"
	@echo ""
	@echo "Docker Targets (前端+后端 统一镜像):"
	@echo "  docker-build  - 一键构建完整Docker镜像 (自动编译前端和后端)"
	@echo "  docker-run    - 运行完整Docker容器 (前端页面+后端API同时启动)"
	@echo "  docker-stop   - 停止并删除Docker容器"
	@echo "  docker-clean  - 清理: 停止容器+删除容器+删除镜像"
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
	@echo "  make all              # 本地启动所有组件 (开发调试)"
	@echo "  make server           # 本地只启动后端"
	@echo "  make web              # 本地只启动前端"
	@echo "  make stop             # 停止所有进程"
	@echo "  make prepare          # 准备数据库环境"
	@echo ""
	@echo "  make build            # 编译所有产物 (前端+后端+客户端)"
	@echo "  make build-web        # 只编译前端"
	@echo "  make build-server     # 只编译后端"
	@echo ""
	@echo "  make docker-build     # 构建Docker镜像 (自动先编译所有产物)"
	@echo "  make docker-run       # 运行Docker容器"
	@echo "  make docker-stop      # 停止Docker容器"
	@echo "  make docker-clean     # 停止+删除容器 + 删除镜像"



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
# 构建前端产物
build-web:
	@echo "构建前端..."
	@cd ./chprobe_web && npm run build
	@echo "✅ 前端构建完成 -> chprobe_web/dist/"

# 构建所有产物 (前端 + 后端 + 客户端)
build:
	@echo "======================================"
	@echo "  编译所有产物"
	@echo "======================================"
	@make build-web
	@make release
	@echo ""
	@echo "✅ 所有产物编译完成:"
	@echo "   前端: chprobe_web/dist/"
	@echo "   后端: chprobe_server/bin/chprobe_server"
	@echo "   客户端: chprobe_client/bin/chprobe_client"
	@echo ""

docker-build:
	@echo "======================================"
	@echo "  构建 ChProbe 统一镜像 (前端+后端)"
	@echo "======================================"
	@echo "→ 第一步: 本地编译所有产物..."
	@make build
	@echo ""
	@echo "→ 第二步: 构建Docker镜像..."
	@sudo docker build -t chprobe -f Dockerfile .
	@echo ""
	@echo "✅ Docker镜像构建完成: chprobe:latest"
	@echo "   镜像包含: 前端页面 + 后端服务 + Nginx"
	@echo ""

# Run unified Docker container
docker-run:
	@make docker-build
	@echo "======================================"
	@echo "  启动 ChProbe 统一容器"
	@echo "======================================"
	@# Stop and remove existing container if it exists
	@sudo docker stop chprobe 2>/dev/null || true
	@sudo docker rm chprobe 2>/dev/null || true
	@# Run new container (使用host网络模式，直接访问宿主机MySQL)
	@sudo docker run --name chprobe \
		--restart always \
		--network host \
		-v chprobe_data:/app/data \
		-d chprobe
	@echo ""
	@echo "✅ 容器启动成功！"
	@echo "   前端页面: http://localhost"
	@echo "   API接口: http://localhost/api"
	@echo "   Agent端口: 32000"
	@echo ""

# Stop and remove unified Docker container
docker-stop:
	@echo "停止并删除 ChProbe 容器..."
	@sudo docker stop chprobe 2>/dev/null || true
	@sudo docker rm chprobe 2>/dev/null || true
	@echo "✅ 容器已停止并删除！"

# Clean Docker: stop + remove container + remove image
docker-clean:
	@echo "======================================"
	@echo "  清理 ChProbe Docker 环境"
	@echo "======================================"
	@echo "→ 停止并删除容器..."
	@sudo docker stop chprobe 2>/dev/null || true
	@sudo docker rm chprobe 2>/dev/null || true
	@echo "→ 删除镜像..."
	@sudo docker rmi chprobe:latest 2>/dev/null || true
	@echo ""
	@echo "✅ Docker 清理完成！"
	@echo "   容器: 已停止并删除"
	@echo "   镜像: 已删除"
	@echo ""

