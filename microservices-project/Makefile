.PHONY: help init gen-user gen-gateway build-user build-gateway run-user run-gateway run stop clean test

# 默认目标
help:
	@echo "Make 命令列表:"
	@echo "  make init          - 初始化项目（安装依赖）"
	@echo "  make gen           - 生成所有代码"
	@echo "  make gen-user      - 生成用户服务代码"
	@echo "  make gen-gateway   - 生成网关代码"
	@echo "  make build         - 编译所有服务"
	@echo "  make build-user    - 编译用户服务"
	@echo "  make build-gateway - 编译网关服务"
	@echo "  make run           - 运行所有服务"
	@echo "  make run-user      - 运行用户服务"
	@echo "  make run-gateway   - 运行网关服务"
	@echo "  make stop          - 停止所有服务"
	@echo "  make clean         - 清理编译文件"
	@echo "  make test          - 运行测试"
	@echo "  make docker-up     - 启动 Docker 服务"
	@echo "  make docker-down   - 停止 Docker 服务"

# 初始化项目
init:
	@echo "==> 初始化项目..."
	go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
	go install github.com/cloudwego/hertz/cmd/hz@latest
	go install github.com/cloudwego/thriftgo@latest
	@echo "==> 安装依赖..."
	cd user-service && go mod download
	cd api-gateway && go mod download

# 生成代码
gen: gen-user gen-gateway

gen-user:
	@echo "==> 生成用户服务代码..."
	cd user-service && kitex -module github.com/example/microservices-project/user-service -service userservice ../idl/user.thrift

gen-gateway:
	@echo "==> 生成网关代码..."
	@echo "网关使用 user-service 的 kitex_gen，无需单独生成"

# 编译
build: build-user build-gateway

build-user:
	@echo "==> 编译用户服务..."
	cd user-service && go build -o bin/user-service main.go

build-gateway:
	@echo "==> 编译网关服务..."
	cd api-gateway && go build -o bin/api-gateway main.go

# 运行
run:
	@echo "==> 启动所有服务..."
	@make run-user &
	@sleep 3
	@make run-gateway &
	@echo "==> 所有服务已启动"

run-user:
	@echo "==> 启动用户服务..."
	cd user-service && go run main.go

run-gateway:
	@echo "==> 启动网关服务..."
	cd api-gateway && go run main.go

# Docker
docker-up:
	@echo "==> 启动 Docker 服务（MySQL + Etcd + Redis）..."
	docker-compose up -d
	@echo "==> 等待服务就绪..."
	sleep 10

docker-down:
	@echo "==> 停止 Docker 服务..."
	docker-compose down

# 停止
stop:
	@echo "==> 停止所有服务..."
	@pkill -f "user-service" || true
	@pkill -f "api-gateway" || true
	@echo "==> 所有服务已停止"

# 清理
clean:
	@echo "==> 清理编译文件..."
	rm -rf user-service/bin
	rm -rf api-gateway/bin
	rm -rf user-service/kitex_gen
	rm -rf api-gateway/kitex_gen
	@echo "==> 清理完成"

# 测试
test:
	@echo "==> 运行测试..."
	cd user-service && go test ./...
	cd api-gateway && go test ./...

# 完整流程
all: docker-up gen build run

.DEFAULT_GOAL := help
