.PHONY: help build run-kitex run-hertz test clean docker-build docker-up docker-down gen-code

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

gen-code: ## 从 IDL 生成代码
	@echo "生成 Kitex 代码..."
	cd kitex_service && kitex -module github.com/example/hertz-kitex-demo -service user_service ../idl/user.thrift

build: ## 编译服务
	@echo "编译 Kitex 服务..."
	go build -o bin/kitex_service ./kitex_service
	@echo "编译 Hertz 服务..."
	go build -o bin/hertz_service ./hertz_service

run-kitex: ## 运行 Kitex RPC 服务
	@echo "启动 Kitex RPC 服务..."
	go run ./kitex_service/main.go ./kitex_service/handler.go

run-hertz: ## 运行 Hertz HTTP 服务
	@echo "启动 Hertz HTTP 服务..."
	go run ./hertz_service/main.go

test: ## 运行测试
	go test -v ./...

clean: ## 清理构建文件
	rm -rf bin/
	go clean

docker-build: ## 构建 Docker 镜像
	docker-compose build

docker-up: ## 启动 Docker 服务
	docker-compose up -d

docker-down: ## 停止 Docker 服务
	docker-compose down

docker-logs: ## 查看 Docker 日志
	docker-compose logs -f

tidy: ## 整理依赖
	go mod tidy

deps: ## 下载依赖
	go mod download
