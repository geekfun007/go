#!/bin/bash

# 启动服务脚本

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "================================"
echo "启动微服务项目"
echo "================================"

# 检查 Docker 服务
echo ""
echo "1. 检查 Docker 服务..."
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}错误: Docker 未运行${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Docker 正在运行${NC}"

# 启动基础服务
echo ""
echo "2. 启动基础服务 (MySQL + Etcd + Redis)..."
docker-compose up -d

echo "   等待服务就绪..."
sleep 10

# 检查服务状态
echo ""
echo "3. 检查服务状态..."
if docker ps | grep -q microservices-mysql; then
    echo -e "${GREEN}✓ MySQL 运行中${NC}"
else
    echo -e "${RED}✗ MySQL 未运行${NC}"
fi

if docker ps | grep -q microservices-etcd; then
    echo -e "${GREEN}✓ Etcd 运行中${NC}"
else
    echo -e "${RED}✗ Etcd 未运行${NC}"
fi

if docker ps | grep -q microservices-redis; then
    echo -e "${GREEN}✓ Redis 运行中${NC}"
else
    echo -e "${RED}✗ Redis 未运行${NC}"
fi

# 启动用户服务
echo ""
echo "4. 启动用户服务..."
cd user-service
nohup go run main.go > ../logs/user-service.log 2>&1 &
USER_SERVICE_PID=$!
echo $USER_SERVICE_PID > ../logs/user-service.pid
echo -e "${GREEN}✓ 用户服务已启动 (PID: $USER_SERVICE_PID)${NC}"
cd ..

sleep 3

# 启动 API 网关
echo ""
echo "5. 启动 API 网关..."
cd api-gateway
nohup go run main.go > ../logs/api-gateway.log 2>&1 &
GATEWAY_PID=$!
echo $GATEWAY_PID > ../logs/api-gateway.pid
echo -e "${GREEN}✓ API 网关已启动 (PID: $GATEWAY_PID)${NC}"
cd ..

echo ""
echo "================================"
echo -e "${GREEN}所有服务启动完成！${NC}"
echo "================================"
echo ""
echo "服务地址："
echo "  - API 网关: http://localhost:8080"
echo "  - 用户服务: localhost:8888"
echo "  - MySQL: localhost:3306"
echo "  - Etcd: localhost:2379"
echo "  - Redis: localhost:6379"
echo ""
echo "日志文件："
echo "  - 用户服务: logs/user-service.log"
echo "  - API 网关: logs/api-gateway.log"
echo ""
echo "测试 API："
echo "  curl http://localhost:8080/health"
echo ""
