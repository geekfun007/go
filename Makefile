.PHONY: help build run-kitex run-hertz test clean docker-build docker-up docker-down gen-kitex gen-hertz update-hertz install-tools

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

install-tools: ## 安装开发工具
	@echo "安装 Kitex CLI..."
	go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
	@echo "安装 Hz CLI..."
	go install github.com/cloudwego/hertz/cmd/hz@latest
	@echo "安装 Thriftgo..."
	go install github.com/cloudwego/thriftgo@latest
	@echo "✅ 工具安装完成"

gen-kitex: ## 从 IDL 生成 Kitex RPC 服务代码
	@echo "生成 Kitex RPC 服务代码..."
	kitex -module github.com/example/hertz-kitex-demo -service user_service ./idl/user.thrift
	@echo "✅ Kitex 代码生成完成"

gen-hertz: ## 从 IDL 生成 Hertz HTTP 服务代码（新项目）
	@echo "生成 Hertz HTTP 服务代码..."
	cd hertz_service && hz new -module github.com/example/hertz-kitex-demo -idl ../idl/user.thrift
	@echo "✅ Hertz 代码生成完成"

update-hertz: ## 更新 Hertz HTTP 服务代码
	@echo "更新 Hertz HTTP 服务代码..."
	cd hertz_service && hz update -idl ../idl/user.thrift
	@echo "✅ Hertz 代码更新完成"

build: ## 编译服务
	@echo "编译 Kitex 服务..."
	cd kitex_service && go build -o ../bin/kitex_service .
	@echo "编译 Hertz 服务..."
	cd hertz_service && go build -o ../bin/hertz_service .
	@echo "✅ 编译完成"

run-kitex: ## 运行 Kitex RPC 服务
	@echo "启动 Kitex RPC 服务..."
	cd kitex_service && go run .

run-hertz: ## 运行 Hertz HTTP 服务
	@echo "启动 Hertz HTTP 服务..."
	cd hertz_service && go run .

test: ## 运行测试
	go test -v ./...

clean: ## 清理构建文件
	rm -rf bin/
	go clean

docker-build: ## 构建 Docker 镜像
	docker-compose build

docker-up: ## 启动 Docker 服务（含 Etcd）
	docker-compose up -d

docker-down: ## 停止 Docker 服务
	docker-compose down

docker-logs: ## 查看 Docker 日志
	docker-compose logs -f

tidy: ## 整理依赖
	go mod tidy
	cd kitex_service && go mod tidy
	cd hertz_service && go mod tidy

deps: ## 下载依赖
	go mod download
