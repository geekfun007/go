#!/bin/bash

# 停止服务脚本

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "================================"
echo "停止微服务项目"
echo "================================"

echo ""
echo "1. 停止用户服务..."
if [ -f logs/user-service.pid ]; then
    PID=$(cat logs/user-service.pid)
    if ps -p $PID > /dev/null; then
        kill $PID
        echo -e "${GREEN}✓ 用户服务已停止${NC}"
    else
        echo -e "${YELLOW}用户服务未运行${NC}"
    fi
    rm -f logs/user-service.pid
else
    # 尝试通过进程名停止
    pkill -f "user-service" || echo -e "${YELLOW}用户服务未运行${NC}"
fi

echo ""
echo "2. 停止 API 网关..."
if [ -f logs/api-gateway.pid ]; then
    PID=$(cat logs/api-gateway.pid)
    if ps -p $PID > /dev/null; then
        kill $PID
        echo -e "${GREEN}✓ API 网关已停止${NC}"
    else
        echo -e "${YELLOW}API 网关未运行${NC}"
    fi
    rm -f logs/api-gateway.pid
else
    pkill -f "api-gateway" || echo -e "${YELLOW}API 网关未运行${NC}"
fi

echo ""
echo "3. 停止 Docker 服务..."
docker-compose down
echo -e "${GREEN}✓ Docker 服务已停止${NC}"

echo ""
echo "================================"
echo -e "${GREEN}所有服务已停止${NC}"
echo "================================"
echo ""
