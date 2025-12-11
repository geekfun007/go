#!/bin/bash

# 代码生成脚本

set -e

echo "================================"
echo "开始生成代码..."
echo "================================"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查工具是否安装
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${YELLOW}警告: $1 未安装，正在安装...${NC}"
        go install $2
    else
        echo -e "${GREEN}✓ $1 已安装${NC}"
    fi
}

echo ""
echo "1. 检查必要工具..."
check_tool "kitex" "github.com/cloudwego/kitex/tool/cmd/kitex@latest"
check_tool "thriftgo" "github.com/cloudwego/thriftgo@latest"

echo ""
echo "2. 生成用户服务代码..."
cd user-service
kitex -module github.com/example/microservices-project/user-service \
      -service userservice \
      ../idl/user.thrift
echo -e "${GREEN}✓ 用户服务代码生成完成${NC}"

cd ..

echo ""
echo "3. 复制 kitex_gen 到 api-gateway..."
rm -rf api-gateway/kitex_gen
cp -r user-service/kitex_gen api-gateway/
echo -e "${GREEN}✓ API 网关代码准备完成${NC}"

echo ""
echo "4. 下载依赖..."
cd user-service && go mod tidy && cd ..
cd api-gateway && go mod tidy && cd ..
echo -e "${GREEN}✓ 依赖下载完成${NC}"

echo ""
echo "================================"
echo -e "${GREEN}代码生成完成！${NC}"
echo "================================"
echo ""
echo "下一步："
echo "  1. 启动基础服务: docker-compose up -d"
echo "  2. 运行用户服务: cd user-service && go run main.go"
echo "  3. 运行 API 网关: cd api-gateway && go run main.go"
echo ""
