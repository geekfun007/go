# API 测试指南

本文档提供完整的 API 测试示例和测试脚本。

## 目录
- [环境准备](#环境准备)
- [API 测试](#api-测试)
- [测试脚本](#测试脚本)
- [Postman 集合](#postman-集合)
- [压力测试](#压力测试)

---

## 环境准备

### 1. 启动服务

```bash
# 启动基础服务
docker-compose up -d

# 等待服务就绪
sleep 10

# 启动用户服务
cd user-service && go run main.go &

# 启动 API 网关
cd api-gateway && go run main.go &
```

### 2. 验证服务状态

```bash
# 检查健康状态
curl http://localhost:8080/health

# 预期响应
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "healthy",
    "services": {
      "api-gateway": "ok",
      "user-service": "ok"
    },
    "timestamp": "2025-12-10T10:00:00Z"
  }
}
```

---

## API 测试

### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "password123",
    "phone": "13800138000",
    "age": 25
  }'
```

**预期响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "testuser@example.com",
      "phone": "13800138000",
      "age": 25,
      "created_at": "2025-12-10T10:00:00Z",
      "updated_at": "2025-12-10T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "password123"
  }'
```

**预期响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "testuser@example.com"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. 获取用户信息

```bash
# 公开接口（无需 token）
curl -X GET http://localhost:8080/api/v1/users/1
```

**预期响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "testuser@example.com",
    "phone": "13800138000",
    "age": 25
  }
}
```

### 4. 获取当前用户信息（需要认证）

```bash
TOKEN="your_jwt_token_here"

curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN"
```

### 5. 更新用户信息

```bash
TOKEN="your_jwt_token_here"

curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "updated_username",
    "age": 26
  }'
```

**预期响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "updated_username",
    "email": "testuser@example.com",
    "age": 26
  }
}
```

### 6. 列出用户（分页）

```bash
# 基本列表
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10"

# 带搜索关键字
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10&keyword=test"
```

**预期响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "testuser",
        "email": "testuser@example.com"
      }
    ],
    "page": {
      "page": 1,
      "page_size": 10,
      "total": 1
    }
  }
}
```

### 7. 删除用户

```bash
TOKEN="your_jwt_token_here"

curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

**预期响应：**
```json
{
  "code": 0,
  "message": "success"
}
```

---

## 测试脚本

### 自动化测试脚本

**test_api.sh**
```bash
#!/bin/bash

# API 自动化测试脚本

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 测试计数器
PASSED=0
FAILED=0

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local expected_code=$5
    local token=$6

    echo "Testing: $name"
    
    if [ -z "$token" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "${BASE_URL}${endpoint}" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)
    
    if [ "$http_code" -eq "$expected_code" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((PASSED++))
        echo "$body" | jq '.'
    else
        echo -e "${RED}✗ FAILED (Expected: $expected_code, Got: $http_code)${NC}"
        ((FAILED++))
        echo "$body"
    fi
    echo ""
}

# 开始测试
echo "================================"
echo "API 自动化测试"
echo "================================"
echo ""

# 1. 健康检查
test_api "Health Check" "GET" "/health" "" 200

# 2. 注册用户
TIMESTAMP=$(date +%s)
REGISTER_DATA="{\"username\":\"testuser${TIMESTAMP}\",\"email\":\"test${TIMESTAMP}@example.com\",\"password\":\"password123\",\"age\":25}"
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REGISTER_DATA")

echo "Testing: Register User"
echo "$REGISTER_RESPONSE" | jq '.'
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.token')
USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user.id')

if [ ! -z "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ PASSED${NC}"
    ((PASSED++))
else
    echo -e "${RED}✗ FAILED${NC}"
    ((FAILED++))
fi
echo ""

# 3. 登录
LOGIN_DATA="{\"email\":\"test${TIMESTAMP}@example.com\",\"password\":\"password123\"}"
test_api "Login User" "POST" "/api/v1/auth/login" "$LOGIN_DATA" 200

# 4. 获取用户信息
test_api "Get User" "GET" "/api/v1/users/${USER_ID}" "" 200

# 5. 获取当前用户信息（需要 token）
test_api "Get Current User" "GET" "/api/v1/users/me" "" 200 "$TOKEN"

# 6. 更新用户信息
UPDATE_DATA="{\"username\":\"updated${TIMESTAMP}\",\"age\":26}"
test_api "Update User" "PUT" "/api/v1/users/${USER_ID}" "$UPDATE_DATA" 200 "$TOKEN"

# 7. 列出用户
test_api "List Users" "GET" "/api/v1/users?page=1&page_size=10" "" 200

# 8. 删除用户
test_api "Delete User" "DELETE" "/api/v1/users/${USER_ID}" "" 200 "$TOKEN"

# 测试结果统计
echo "================================"
echo "测试结果"
echo "================================"
echo -e "通过: ${GREEN}${PASSED}${NC}"
echo -e "失败: ${RED}${FAILED}${NC}"
echo "总计: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}有测试失败！${NC}"
    exit 1
fi
```

### 运行测试

```bash
chmod +x test_api.sh
./test_api.sh
```

---

## Postman 集合

### 导入 Postman 集合

**microservices.postman_collection.json**
```json
{
  "info": {
    "name": "Microservices API",
    "description": "GORM + Hertz + Kitex 微服务 API 测试集合",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/health",
          "host": ["{{base_url}}"],
          "path": ["health"]
        }
      }
    },
    {
      "name": "Register",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "var jsonData = pm.response.json();",
              "if (jsonData.data && jsonData.data.token) {",
              "    pm.environment.set('token', jsonData.data.token);",
              "    pm.environment.set('user_id', jsonData.data.user.id);",
              "}"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"username\": \"testuser\",\n  \"email\": \"test@example.com\",\n  \"password\": \"password123\",\n  \"phone\": \"13800138000\",\n  \"age\": 25\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/v1/auth/register",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "register"]
        }
      }
    },
    {
      "name": "Login",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "var jsonData = pm.response.json();",
              "if (jsonData.data && jsonData.data.token) {",
              "    pm.environment.set('token', jsonData.data.token);",
              "}"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"password123\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/v1/auth/login",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "auth", "login"]
        }
      }
    },
    {
      "name": "Get User",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/v1/users/{{user_id}}",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "users", "{{user_id}}"]
        }
      }
    },
    {
      "name": "Get Current User",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": {
          "raw": "{{base_url}}/api/v1/users/me",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "users", "me"]
        }
      }
    },
    {
      "name": "Update User",
      "request": {
        "method": "PUT",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          },
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"username\": \"updated_username\",\n  \"age\": 26\n}"
        },
        "url": {
          "raw": "{{base_url}}/api/v1/users/{{user_id}}",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "users", "{{user_id}}"]
        }
      }
    },
    {
      "name": "List Users",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "{{base_url}}/api/v1/users?page=1&page_size=10",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "users"],
          "query": [
            {
              "key": "page",
              "value": "1"
            },
            {
              "key": "page_size",
              "value": "10"
            }
          ]
        }
      }
    },
    {
      "name": "Delete User",
      "request": {
        "method": "DELETE",
        "header": [
          {
            "key": "Authorization",
            "value": "Bearer {{token}}"
          }
        ],
        "url": {
          "raw": "{{base_url}}/api/v1/users/{{user_id}}",
          "host": ["{{base_url}}"],
          "path": ["api", "v1", "users", "{{user_id}}"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080"
    }
  ]
}
```

---

## 压力测试

### 使用 Apache Bench (ab)

```bash
# 安装 ab
sudo apt-get install apache2-utils

# 注册接口压测
ab -n 1000 -c 10 -p register.json -T application/json \
  http://localhost:8080/api/v1/auth/register

# 登录接口压测
ab -n 1000 -c 10 -p login.json -T application/json \
  http://localhost:8080/api/v1/auth/login
```

**register.json**
```json
{
  "username": "loadtest",
  "email": "loadtest@example.com",
  "password": "password123",
  "age": 25
}
```

**login.json**
```json
{
  "email": "loadtest@example.com",
  "password": "password123"
}
```

### 使用 wrk

```bash
# 安装 wrk
sudo apt-get install wrk

# 运行压测
wrk -t12 -c400 -d30s http://localhost:8080/health

# 带 POST 数据
wrk -t12 -c400 -d30s -s post.lua http://localhost:8080/api/v1/auth/login
```

**post.lua**
```lua
wrk.method = "POST"
wrk.body   = '{"email":"test@example.com","password":"password123"}'
wrk.headers["Content-Type"] = "application/json"
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |
| 40001 | 邮箱已存在 |
| 40002 | 用户名已存在 |
| 40003 | 用户名已被占用 |
| 40004 | 邮箱已被占用 |
| 40101 | 密码错误 |
| 40401 | 用户不存在（登录）|
| 40402 | 用户不存在（查询）|
| 40403 | 用户不存在（更新）|
| 40404 | 用户不存在（删除）|

---

## 常见问题

### Q: 获取 token 后如何使用？

A: 在请求头中添加：
```
Authorization: Bearer {your_token}
```

### Q: 如何批量测试？

A: 使用提供的 `test_api.sh` 脚本，或使用 Postman 的 Collection Runner。

### Q: 如何查看详细日志？

A: 查看日志文件：
- 用户服务：`logs/user-service.log`
- API 网关：`logs/api-gateway.log`

### Q: 压测时注意什么？

A: 
1. 注意数据库连接池配置
2. 调整系统文件描述符限制
3. 监控服务器资源使用情况
4. 使用不同的测试数据避免唯一键冲突
