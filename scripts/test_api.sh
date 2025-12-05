#!/bin/bash

# API 测试脚本

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "测试 Hertz + Kitex + MySQL 服务"
echo "=========================================="
echo ""

# 1. 健康检查
echo "1. 健康检查..."
curl -s $BASE_URL/ping
echo -e "\n"

# 2. 创建用户
echo "2. 创建用户..."
USER1=$(curl -s -X POST $BASE_URL/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "phone": "13800138001"
  }')
echo $USER1 | jq '.'
USER1_ID=$(echo $USER1 | jq -r '.data.id')
echo ""

# 3. 创建第二个用户
echo "3. 创建第二个用户..."
USER2=$(curl -s -X POST $BASE_URL/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "email": "bob@example.com",
    "phone": "13800138002"
  }')
echo $USER2 | jq '.'
USER2_ID=$(echo $USER2 | jq -r '.data.id')
echo ""

# 4. 获取用户
echo "4. 获取用户 (ID: $USER1_ID)..."
curl -s $BASE_URL/api/v1/users/$USER1_ID | jq '.'
echo ""

# 5. 获取用户列表
echo "5. 获取用户列表..."
curl -s "$BASE_URL/api/v1/users?page=1&page_size=10" | jq '.'
echo ""

# 6. 更新用户
echo "6. 更新用户 (ID: $USER1_ID)..."
curl -s -X PUT $BASE_URL/api/v1/users/$USER1_ID \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice_updated",
    "email": "alice_new@example.com"
  }' | jq '.'
echo ""

# 7. 验证更新
echo "7. 验证更新后的用户信息..."
curl -s $BASE_URL/api/v1/users/$USER1_ID | jq '.'
echo ""

# 8. 删除用户
echo "8. 删除用户 (ID: $USER2_ID)..."
curl -s -X DELETE $BASE_URL/api/v1/users/$USER2_ID | jq '.'
echo ""

# 9. 验证删除
echo "9. 验证删除（应该返回 404）..."
curl -s $BASE_URL/api/v1/users/$USER2_ID | jq '.'
echo ""

# 10. 最终用户列表
echo "10. 最终用户列表..."
curl -s "$BASE_URL/api/v1/users?page=1&page_size=10" | jq '.'
echo ""

echo "=========================================="
echo "测试完成！"
echo "=========================================="
