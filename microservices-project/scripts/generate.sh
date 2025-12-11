#!/bin/bash

# 代码生成脚本（更新版：支持 Hertz 和 Kitex）

set -e

echo "================================"
echo "开始生成代码..."
echo "================================"

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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
echo -e "${BLUE}步骤 1: 检查必要工具...${NC}"
check_tool "kitex" "github.com/cloudwego/kitex/tool/cmd/kitex@latest"
check_tool "hz" "github.com/cloudwego/hertz/cmd/hz@latest"
check_tool "thriftgo" "github.com/cloudwego/thriftgo@latest"

echo ""
echo -e "${BLUE}步骤 2: 生成用户服务代码 (Kitex)...${NC}"
cd user-service

# 清理旧的生成代码
rm -rf kitex_gen

# 生成 Kitex 服务端代码
kitex -module github.com/example/microservices-project/user-service \
      -service userservice \
      -use github.com/example/microservices-project/user-service/kitex_gen \
      ../idl/user.thrift

echo -e "${GREEN}✓ 用户服务代码生成完成${NC}"
cd ..

echo ""
echo -e "${BLUE}步骤 3: 生成 API 网关代码 (Hertz)...${NC}"
cd api-gateway

# 清理旧的生成代码
rm -rf biz/handler
rm -rf biz/router
rm -rf biz/model
rm -rf hertz_gen

# 使用 hz 生成 Hertz 代码
# 注意：hz 会生成 handler、router、model 等
hz new -module github.com/example/microservices-project/api-gateway \
   -idl ../idl/user.thrift \
   -handler_dir biz/handler \
   -model_dir biz/model \
   -router_dir biz/router \
   -force

echo -e "${GREEN}✓ API 网关代码生成完成${NC}"
cd ..

echo ""
echo -e "${BLUE}步骤 4: 复制 kitex_gen 到 api-gateway...${NC}"
# API 网关需要调用用户服务，所以需要 kitex_gen
cp -r user-service/kitex_gen api-gateway/
echo -e "${GREEN}✓ kitex_gen 复制完成${NC}"

echo ""
echo -e "${BLUE}步骤 5: 下载依赖...${NC}"
echo "  - 用户服务..."
cd user-service && go mod tidy && cd ..
echo "  - API 网关..."
cd api-gateway && go mod tidy && cd ..
echo -e "${GREEN}✓ 依赖下载完成${NC}"

echo ""
echo "================================"
echo -e "${GREEN}代码生成完成！${NC}"
echo "================================"
echo ""
echo -e "${YELLOW}重要说明：${NC}"
echo "  1. hz 已自动生成 Handler 和 Router 代码"
echo "  2. 生成的 handler 在: api-gateway/biz/handler/"
echo "  3. 生成的 router 在: api-gateway/biz/router/"
echo "  4. 你需要在生成的 handler 中实现业务逻辑"
echo ""
echo -e "${BLUE}下一步：${NC}"
echo "  1. 修改 api-gateway/biz/handler 中的实现"
echo "  2. 添加中间件: api-gateway/biz/middleware/"
echo "  3. 修改 main.go 注册中间件"
echo "  4. 启动服务: docker-compose up -d"
echo "  5. 运行用户服务: cd user-service && go run main.go"
echo "  6. 运行 API 网关: cd api-gateway && go run main.go"
echo ""
