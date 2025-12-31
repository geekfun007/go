# Docker & Docker Compose 详解 / Docker & Docker Compose Guide

## 1. Dockerfile 基础 / Dockerfile Basics

### 1.1 基本 Go 应用 Dockerfile / Basic Go Application Dockerfile

```dockerfile
# ========================================
# 简单 Dockerfile / Simple Dockerfile
# ========================================

FROM golang:1.21-alpine

WORKDIR /app

# 复制 go mod 文件 / Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码 / Copy source code
COPY . .

# 构建 / Build
RUN go build -o main .

# 运行 / Run
EXPOSE 8080
CMD ["./main"]
```

### 1.2 多阶段构建 (推荐) / Multi-stage Build (Recommended)

```dockerfile
# ========================================
# 多阶段构建 - 生产环境推荐
# Multi-stage build - Recommended for production
# ========================================

# 阶段1: 构建 / Stage 1: Build
FROM golang:1.21-alpine AS builder

# 安装必要工具 / Install necessary tools
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# 复制依赖文件 / Copy dependency files
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 复制源代码 / Copy source code
COPY . .

# 构建静态二进制 / Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /app/server ./cmd/server

# 阶段2: 运行 / Stage 2: Run
FROM scratch

# 从 builder 复制证书和时区 / Copy certs and timezone from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 复制二进制 / Copy binary
COPY --from=builder /app/server /server

# 复制配置文件 (如果需要) / Copy config files (if needed)
COPY --from=builder /app/config /config

# 设置时区 / Set timezone
ENV TZ=Asia/Shanghai

# 非 root 用户 / Non-root user
USER 1000:1000

EXPOSE 8080

ENTRYPOINT ["/server"]
```

### 1.3 带 Alpine 的多阶段构建 / Multi-stage with Alpine

```dockerfile
# ========================================
# 使用 Alpine 作为运行镜像
# Using Alpine as runtime image
# ========================================

# 构建阶段 / Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO 启用时的构建 / Build with CGO enabled
RUN go build -o server ./cmd/server

# 运行阶段 / Runtime stage
FROM alpine:3.19

# 安装基本工具和证书 / Install basic tools and certs
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 复制二进制 / Copy binary
COPY --from=builder /app/server .

# 复制静态资源 / Copy static assets
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# 创建非 root 用户 / Create non-root user
RUN adduser -D -g '' appuser
USER appuser

ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./server"]
```

### 1.4 开发环境 Dockerfile / Development Dockerfile

```dockerfile
# ========================================
# 开发环境 Dockerfile
# Development Dockerfile
# ========================================

FROM golang:1.21

# 安装开发工具 / Install development tools
RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# 复制依赖 / Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# 源代码通过 volume 挂载 / Source code mounted via volume

EXPOSE 8080
EXPOSE 2345

# 使用 air 进行热重载 / Use air for hot reload
CMD ["air", "-c", ".air.toml"]
```

## 2. Docker Compose 配置 / Docker Compose Configuration

### 2.1 基本配置 / Basic Configuration

```yaml
# docker-compose.yml
version: '3.8'

services:
  # ========================================
  # Go 应用服务 / Go Application Service
  # ========================================
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: go-app
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=password
      - DB_NAME=myapp
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app-network
    restart: unless-stopped

  # ========================================
  # MySQL 服务 / MySQL Service
  # ========================================
  mysql:
    image: mysql:8.0
    container_name: mysql
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: myapp
      MYSQL_USER: appuser
      MYSQL_PASSWORD: apppassword
    volumes:
      - mysql-data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network
    restart: unless-stopped

  # ========================================
  # Redis 服务 / Redis Service
  # ========================================
  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    command: redis-server --appendonly yes --requirepass redispassword
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "redispassword", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network
    restart: unless-stopped

networks:
  app-network:
    driver: bridge

volumes:
  mysql-data:
  redis-data:
```

### 2.2 开发环境配置 / Development Configuration

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: go-app-dev
    ports:
      - "8080:8080"
      - "2345:2345"  # Delve debugger
    environment:
      - APP_ENV=development
      - DB_HOST=mysql
      - DB_PORT=3306
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - .:/app  # 挂载源代码 / Mount source code
      - go-modules:/go/pkg/mod  # 缓存依赖 / Cache dependencies
    depends_on:
      - mysql
      - redis
    networks:
      - dev-network

  mysql:
    image: mysql:8.0
    container_name: mysql-dev
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: devpassword
      MYSQL_DATABASE: devdb
    volumes:
      - mysql-dev-data:/var/lib/mysql
    networks:
      - dev-network

  redis:
    image: redis:7-alpine
    container_name: redis-dev
    ports:
      - "6379:6379"
    networks:
      - dev-network

  # ========================================
  # 开发工具 / Development Tools
  # ========================================
  
  # Adminer - 数据库管理 / Database management
  adminer:
    image: adminer
    container_name: adminer
    ports:
      - "8081:8080"
    networks:
      - dev-network

  # Redis Commander - Redis 管理 / Redis management
  redis-commander:
    image: rediscommander/redis-commander:latest
    container_name: redis-commander
    ports:
      - "8082:8081"
    environment:
      - REDIS_HOSTS=local:redis:6379
    networks:
      - dev-network

networks:
  dev-network:
    driver: bridge

volumes:
  mysql-dev-data:
  go-modules:
```

### 2.3 完整微服务配置 / Complete Microservices Configuration

```yaml
# docker-compose.microservices.yml
version: '3.8'

services:
  # ========================================
  # API Gateway / API 网关
  # ========================================
  gateway:
    build:
      context: ./gateway
      dockerfile: Dockerfile
    container_name: gateway
    ports:
      - "80:8080"
    environment:
      - USER_SERVICE_URL=http://user-service:8080
      - ORDER_SERVICE_URL=http://order-service:8080
    depends_on:
      - user-service
      - order-service
    networks:
      - microservices

  # ========================================
  # User Service / 用户服务
  # ========================================
  user-service:
    build:
      context: ./services/user
      dockerfile: Dockerfile
    container_name: user-service
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis
    networks:
      - microservices
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '0.5'
          memory: 256M

  # ========================================
  # Order Service / 订单服务
  # ========================================
  order-service:
    build:
      context: ./services/order
      dockerfile: Dockerfile
    container_name: order-service
    environment:
      - DB_HOST=mysql
      - REDIS_HOST=redis
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - mysql
      - redis
      - kafka
    networks:
      - microservices
    deploy:
      replicas: 2

  # ========================================
  # Infrastructure / 基础设施
  # ========================================
  mysql:
    image: mysql:8.0
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD:-password}
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - microservices
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    container_name: redis
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    networks:
      - microservices

  # ========================================
  # Kafka / 消息队列
  # ========================================
  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
    networks:
      - microservices

  kafka:
    image: confluentinc/cp-kafka:7.4.0
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - microservices

  # ========================================
  # Monitoring / 监控
  # ========================================
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - microservices

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-data:/var/lib/grafana
    networks:
      - microservices

networks:
  microservices:
    driver: bridge

volumes:
  mysql-data:
  redis-data:
  grafana-data:
```

## 3. Docker 常用命令 / Docker Common Commands

```bash
# ========================================
# 镜像操作 / Image Operations
# ========================================

# 构建镜像 / Build image
docker build -t myapp:latest .
docker build -t myapp:v1.0 -f Dockerfile.prod .

# 查看镜像 / List images
docker images

# 删除镜像 / Remove image
docker rmi myapp:latest

# 推送镜像 / Push image
docker push registry.example.com/myapp:latest

# ========================================
# 容器操作 / Container Operations
# ========================================

# 运行容器 / Run container
docker run -d --name myapp -p 8080:8080 myapp:latest

# 带环境变量运行 / Run with environment variables
docker run -d --name myapp \
    -e DB_HOST=localhost \
    -e DB_PORT=3306 \
    -p 8080:8080 \
    myapp:latest

# 挂载卷运行 / Run with volumes
docker run -d --name myapp \
    -v $(pwd)/config:/app/config \
    -v app-data:/app/data \
    -p 8080:8080 \
    myapp:latest

# 查看容器 / List containers
docker ps
docker ps -a

# 进入容器 / Enter container
docker exec -it myapp /bin/sh

# 查看日志 / View logs
docker logs myapp
docker logs -f myapp  # 实时日志
docker logs --tail 100 myapp  # 最后100行

# 停止/启动/重启容器 / Stop/Start/Restart container
docker stop myapp
docker start myapp
docker restart myapp

# 删除容器 / Remove container
docker rm myapp
docker rm -f myapp  # 强制删除运行中的容器

# ========================================
# Docker Compose 操作 / Docker Compose Operations
# ========================================

# 启动服务 / Start services
docker-compose up
docker-compose up -d  # 后台运行
docker-compose up --build  # 重新构建

# 使用特定配置文件 / Use specific config file
docker-compose -f docker-compose.dev.yml up -d

# 查看服务状态 / View service status
docker-compose ps

# 查看日志 / View logs
docker-compose logs
docker-compose logs -f app  # 实时查看特定服务

# 停止服务 / Stop services
docker-compose down
docker-compose down -v  # 同时删除卷

# 重启服务 / Restart services
docker-compose restart
docker-compose restart app

# 扩展服务 / Scale services
docker-compose up -d --scale app=3

# 执行命令 / Execute commands
docker-compose exec app /bin/sh
docker-compose exec mysql mysql -u root -p

# ========================================
# 清理操作 / Cleanup Operations
# ========================================

# 清理未使用资源 / Clean unused resources
docker system prune
docker system prune -a  # 包括未使用的镜像

# 清理卷 / Clean volumes
docker volume prune

# 清理网络 / Clean networks
docker network prune
```

## 4. Go 应用配置示例 / Go Application Configuration Example

```go
// config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
}

type ServerConfig struct {
    Host string
    Port int
    Env  string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

func Load() *Config {
    return &Config{
        Server: ServerConfig{
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
            Port: getEnvAsInt("SERVER_PORT", 8080),
            Env:  getEnv("APP_ENV", "development"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 3306),
            User:     getEnv("DB_USER", "root"),
            Password: getEnv("DB_PASSWORD", ""),
            DBName:   getEnv("DB_NAME", "myapp"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
        },
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}
```

## 5. 健康检查端点 / Health Check Endpoint

```go
// handlers/health.go
package handlers

import (
    "context"
    "database/sql"
    "net/http"
    "time"

    "github.com/redis/go-redis/v9"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services"`
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{db: db, redis: redis}
}

func (h *HealthChecker) Check(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    status := HealthStatus{
        Status:    "healthy",
        Timestamp: time.Now(),
        Services:  make(map[string]string),
    }

    // 检查数据库 / Check database
    if err := h.db.PingContext(ctx); err != nil {
        status.Services["database"] = "unhealthy: " + err.Error()
        status.Status = "unhealthy"
    } else {
        status.Services["database"] = "healthy"
    }

    // 检查 Redis / Check Redis
    if err := h.redis.Ping(ctx).Err(); err != nil {
        status.Services["redis"] = "unhealthy: " + err.Error()
        status.Status = "unhealthy"
    } else {
        status.Services["redis"] = "healthy"
    }

    statusCode := http.StatusOK
    if status.Status == "unhealthy" {
        statusCode = http.StatusServiceUnavailable
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(status)
}

// Liveness - Kubernetes liveness probe
func (h *HealthChecker) Liveness(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

// Readiness - Kubernetes readiness probe
func (h *HealthChecker) Readiness(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    // 检查关键依赖 / Check critical dependencies
    if err := h.db.PingContext(ctx); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Database not ready"))
        return
    }

    if err := h.redis.Ping(ctx).Err(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Redis not ready"))
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Ready"))
}
```

## 6. .dockerignore 文件 / .dockerignore File

```
# .dockerignore

# Git
.git
.gitignore

# IDE
.idea
.vscode
*.swp
*.swo

# 依赖目录 / Dependencies
vendor/

# 测试文件 / Test files
*_test.go
**/*_test.go

# 文档 / Documentation
*.md
docs/

# CI/CD
.github/
.gitlab-ci.yml
Jenkinsfile

# Docker
Dockerfile*
docker-compose*.yml
.docker/

# 本地配置 / Local config
.env
.env.*
config.local.yaml

# 构建产物 / Build artifacts
bin/
dist/
*.exe

# 日志 / Logs
*.log
logs/

# 临时文件 / Temp files
tmp/
temp/
```
