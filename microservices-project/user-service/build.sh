#!/bin/bash

# 构建脚本

set -e

echo "==> Building user-service..."

# 创建 bin 目录
mkdir -p bin

# 编译
go build -o bin/user-service main.go

echo "==> Build complete: bin/user-service"

# 运行（可选）
# ./bin/user-service
