# Redis 详解 (GORM 集成) / Redis Guide (GORM Integration)

## 1. 连接与配置 / Connection & Configuration

### 1.1 基础连接 (go-redis) / Basic Connection

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

// 安装 / Installation:
// go get github.com/redis/go-redis/v9

// Redis 客户端配置 / Redis client configuration
type RedisConfig struct {
    Addr         string
    Password     string
    DB           int
    DialTimeout  time.Duration
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    PoolSize     int
    MinIdleConns int
}

// 创建 Redis 客户端 / Create Redis client
func NewRedisClient(cfg *RedisConfig) *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         cfg.Addr,
        Password:     cfg.Password,
        DB:           cfg.DB,
        DialTimeout:  cfg.DialTimeout,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
    })
}

func main() {
    ctx := context.Background()

    // 单机连接 / Standalone Connection
    rdb := NewRedisClient(&RedisConfig{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolSize:     10,
        MinIdleConns: 5,
    })
    defer rdb.Close()

    // 测试连接 / Test connection
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected:", pong)
}
```

### 1.2 GORM + Redis 集成架构 / GORM + Redis Integration

```go
package main

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// 安装 / Installation:
// go get gorm.io/gorm
// go get gorm.io/driver/mysql
// go get github.com/redis/go-redis/v9

// 数据库服务结构 / Database service structure
type DBService struct {
    DB    *gorm.DB
    Redis *redis.Client
}

// 创建数据库服务 / Create database service
func NewDBService(mysqlDSN string, redisAddr string) (*DBService, error) {
    // 连接 MySQL (使用 GORM)
    db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // 配置 MySQL 连接池
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    // 连接 Redis
    rdb := redis.NewClient(&redis.Options{
        Addr:     redisAddr,
        PoolSize: 10,
    })

    return &DBService{
        DB:    db,
        Redis: rdb,
    }, nil
}

// 关闭连接 / Close connections
func (s *DBService) Close() error {
    sqlDB, _ := s.DB.DB()
    sqlDB.Close()
    return s.Redis.Close()
}

// 使用示例 / Usage example
func main() {
    service, err := NewDBService(
        "user:pass@tcp(localhost:3306)/dbname?parseTime=true",
        "localhost:6379",
    )
    if err != nil {
        panic(err)
    }
    defer service.Close()

    // 现在可以同时使用 GORM 和 Redis
    // Now can use both GORM and Redis
}
```

### 1.3 集群和哨兵模式 / Cluster & Sentinel Mode

```go
// 哨兵模式 / Sentinel Mode
func NewSentinelClient() *redis.Client {
    return redis.NewFailoverClient(&redis.FailoverOptions{
        MasterName:    "mymaster",
        SentinelAddrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
        Password:      "",
        DB:            0,
    })
}

// 集群模式 / Cluster Mode
func NewClusterClient() *redis.ClusterClient {
    return redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001",
            "localhost:7002",
        },
        Password:     "",
        PoolSize:     10,
        MinIdleConns: 5,
    })
}
```

## 2. GORM 模型与 Redis 缓存集成 / GORM Models with Redis Caching

### 2.1 缓存层实现 / Cache Layer Implementation

```go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

// 用户模型 (GORM) / User model (GORM)
type User struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Name      string    `gorm:"size:100" json:"name"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// 用户仓库 (带缓存) / User repository (with cache)
type UserRepository struct {
    db       *gorm.DB
    redis    *redis.Client
    cacheTTL time.Duration
}

func NewUserRepository(db *gorm.DB, rdb *redis.Client) *UserRepository {
    return &UserRepository{
        db:       db,
        redis:    rdb,
        cacheTTL: 10 * time.Minute,
    }
}

// 缓存 key 生成 / Cache key generation
func (r *UserRepository) cacheKey(id uint) string {
    return fmt.Sprintf("user:%d", id)
}

// 从缓存获取 / Get from cache
func (r *UserRepository) getFromCache(ctx context.Context, id uint) (*User, error) {
    data, err := r.redis.Get(ctx, r.cacheKey(id)).Bytes()
    if err != nil {
        return nil, err
    }

    var user User
    if err := json.Unmarshal(data, &user); err != nil {
        return nil, err
    }
    return &user, nil
}

// 写入缓存 / Set to cache
func (r *UserRepository) setToCache(ctx context.Context, user *User) error {
    data, err := json.Marshal(user)
    if err != nil {
        return err
    }
    return r.redis.Set(ctx, r.cacheKey(user.ID), data, r.cacheTTL).Err()
}

// 删除缓存 / Delete from cache
func (r *UserRepository) deleteFromCache(ctx context.Context, id uint) error {
    return r.redis.Del(ctx, r.cacheKey(id)).Err()
}

// GetByID - 先查缓存，再查数据库 / Check cache first, then database
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*User, error) {
    // 1. 尝试从缓存获取 / Try to get from cache
    user, err := r.getFromCache(ctx, id)
    if err == nil {
        return user, nil
    }

    // 2. 缓存未命中，从数据库查询 / Cache miss, query from database
    user = &User{}
    if err := r.db.First(user, id).Error; err != nil {
        return nil, err
    }

    // 3. 写入缓存 / Write to cache
    r.setToCache(ctx, user)

    return user, nil
}

// Create - 创建并写入缓存 / Create and cache
func (r *UserRepository) Create(ctx context.Context, user *User) error {
    // 使用 GORM 创建
    if err := r.db.Create(user).Error; err != nil {
        return err
    }

    // 写入缓存
    return r.setToCache(ctx, user)
}

// Update - 更新数据库并刷新缓存 / Update database and refresh cache
func (r *UserRepository) Update(ctx context.Context, user *User) error {
    // 更新数据库
    if err := r.db.Save(user).Error; err != nil {
        return err
    }

    // 刷新缓存
    return r.setToCache(ctx, user)
}

// Delete - 删除数据库记录并清除缓存 / Delete from database and clear cache
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    // 删除缓存
    r.deleteFromCache(ctx, id)

    // 删除数据库记录
    return r.db.Delete(&User{}, id).Error
}
```

### 2.2 使用 singleflight 防止缓存击穿 / Prevent Cache Breakdown with singleflight

```go
import "golang.org/x/sync/singleflight"

type CachedUserRepository struct {
    db       *gorm.DB
    redis    *redis.Client
    cacheTTL time.Duration
    sf       singleflight.Group
}

func (r *CachedUserRepository) GetByID(ctx context.Context, id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 1. 尝试从缓存获取
    data, err := r.redis.Get(ctx, cacheKey).Bytes()
    if err == nil {
        var user User
        json.Unmarshal(data, &user)
        return &user, nil
    }

    // 2. 使用 singleflight 防止缓存击穿
    v, err, _ := r.sf.Do(cacheKey, func() (interface{}, error) {
        // 再次检查缓存 (double check)
        data, err := r.redis.Get(ctx, cacheKey).Bytes()
        if err == nil {
            var user User
            json.Unmarshal(data, &user)
            return &user, nil
        }

        // 从数据库查询
        var user User
        if err := r.db.First(&user, id).Error; err != nil {
            return nil, err
        }

        // 写入缓存
        data, _ := json.Marshal(user)
        r.redis.Set(ctx, cacheKey, data, r.cacheTTL)

        return &user, nil
    })

    if err != nil {
        return nil, err
    }
    return v.(*User), nil
}
```

## 3. 基本数据类型操作 / Basic Data Type Operations

### 3.1 String 字符串

```go
func stringOperations(ctx context.Context, rdb *redis.Client) {
    // SET / GET
    rdb.Set(ctx, "name", "Alice", 10*time.Minute)
    val, _ := rdb.Get(ctx, "name").Result()
    fmt.Println("name:", val)

    // SetNX - 仅当 key 不存在时设置 (用于分布式锁)
    wasSet, _ := rdb.SetNX(ctx, "lock", "1", 30*time.Second).Result()
    fmt.Println("SetNX result:", wasSet)

    // 批量操作
    rdb.MSet(ctx, "key1", "val1", "key2", "val2")
    vals, _ := rdb.MGet(ctx, "key1", "key2").Result()
    fmt.Println("MGet:", vals)

    // 数值操作
    rdb.Set(ctx, "counter", 0, 0)
    rdb.Incr(ctx, "counter")
    rdb.IncrBy(ctx, "counter", 10)
    count, _ := rdb.Get(ctx, "counter").Int64()
    fmt.Println("counter:", count)
}
```

### 3.2 Hash 哈希 - 存储 GORM 模型

```go
func hashWithGORMModel(ctx context.Context, rdb *redis.Client) {
    // 将 GORM 模型存储到 Redis Hash
    user := User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   25,
    }

    // 存储为 Hash
    rdb.HSet(ctx, "user:1", map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "age":   user.Age,
    })

    // 读取
    data, _ := rdb.HGetAll(ctx, "user:1").Result()
    fmt.Println("User:", data)

    // 使用 Scan 读取到结构体
    type UserCache struct {
        ID    int64  `redis:"id"`
        Name  string `redis:"name"`
        Email string `redis:"email"`
        Age   int    `redis:"age"`
    }

    var userCache UserCache
    rdb.HGetAll(ctx, "user:1").Scan(&userCache)
    fmt.Printf("UserCache: %+v\n", userCache)
}
```

### 3.3 List 列表 - 消息队列

```go
func listAsQueue(ctx context.Context, rdb *redis.Client) {
    // 生产者 - 添加任务
    rdb.RPush(ctx, "tasks", "task1", "task2", "task3")

    // 消费者 - 阻塞获取任务
    result, err := rdb.BLPop(ctx, 5*time.Second, "tasks").Result()
    if err == redis.Nil {
        fmt.Println("No task available")
    } else {
        fmt.Println("Processing:", result[1])
    }

    // 获取队列长度
    length, _ := rdb.LLen(ctx, "tasks").Result()
    fmt.Println("Queue length:", length)
}
```

### 3.4 Set 集合 - 标签系统

```go
func setAsTags(ctx context.Context, rdb *redis.Client) {
    // 添加用户标签
    rdb.SAdd(ctx, "user:1:tags", "go", "redis", "docker")
    rdb.SAdd(ctx, "user:2:tags", "redis", "kubernetes", "python")

    // 获取用户标签
    tags, _ := rdb.SMembers(ctx, "user:1:tags").Result()
    fmt.Println("User 1 tags:", tags)

    // 共同标签 (交集)
    commonTags, _ := rdb.SInter(ctx, "user:1:tags", "user:2:tags").Result()
    fmt.Println("Common tags:", commonTags)

    // 所有标签 (并集)
    allTags, _ := rdb.SUnion(ctx, "user:1:tags", "user:2:tags").Result()
    fmt.Println("All tags:", allTags)
}
```

### 3.5 Sorted Set 有序集合 - 排行榜

```go
func sortedSetAsLeaderboard(ctx context.Context, rdb *redis.Client) {
    // 添加玩家分数
    rdb.ZAdd(ctx, "leaderboard",
        redis.Z{Score: 100, Member: "player1"},
        redis.Z{Score: 85, Member: "player2"},
        redis.Z{Score: 92, Member: "player3"},
    )

    // 增加分数
    rdb.ZIncrBy(ctx, "leaderboard", 5, "player2")

    // 获取排名 (降序)
    rank, _ := rdb.ZRevRank(ctx, "leaderboard", "player1").Result()
    fmt.Println("Player1 rank:", rank+1)  // 排名从0开始

    // 获取 Top 3
    top3, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()
    fmt.Println("Top 3:")
    for i, z := range top3 {
        fmt.Printf("  %d. %s: %.0f\n", i+1, z.Member, z.Score)
    }

    // 获取分数范围内的玩家
    players, _ := rdb.ZRangeByScore(ctx, "leaderboard", &redis.ZRangeBy{
        Min: "80",
        Max: "100",
    }).Result()
    fmt.Println("Players with score 80-100:", players)
}
```

## 4. 高级功能 / Advanced Features

### 4.1 Pipeline 管道 - 批量操作

```go
func pipelineWithGORM(ctx context.Context, db *gorm.DB, rdb *redis.Client) {
    // 批量缓存多个用户
    var users []User
    db.Limit(100).Find(&users)

    pipe := rdb.Pipeline()
    for _, user := range users {
        data, _ := json.Marshal(user)
        pipe.Set(ctx, fmt.Sprintf("user:%d", user.ID), data, 10*time.Minute)
    }
    _, err := pipe.Exec(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Cached %d users\n", len(users))

    // 批量读取
    pipe = rdb.Pipeline()
    gets := make([]*redis.StringCmd, len(users))
    for i, user := range users {
        gets[i] = pipe.Get(ctx, fmt.Sprintf("user:%d", user.ID))
    }
    pipe.Exec(ctx)

    for i, get := range gets {
        if get.Err() == nil {
            fmt.Printf("User %d cached\n", users[i].ID)
        }
    }
}
```

### 4.2 事务 (Watch) / Transaction

```go
func transactionExample(ctx context.Context, rdb *redis.Client) {
    // 乐观锁实现转账
    transfer := func(from, to string, amount int64) error {
        return rdb.Watch(ctx, func(tx *redis.Tx) error {
            fromBalance, err := tx.Get(ctx, from).Int64()
            if err != nil && err != redis.Nil {
                return err
            }

            if fromBalance < amount {
                return errors.New("insufficient balance")
            }

            toBalance, _ := tx.Get(ctx, to).Int64()

            _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
                pipe.Set(ctx, from, fromBalance-amount, 0)
                pipe.Set(ctx, to, toBalance+amount, 0)
                return nil
            })
            return err
        }, from, to)
    }

    rdb.Set(ctx, "account:A", 100, 0)
    rdb.Set(ctx, "account:B", 50, 0)

    err := transfer("account:A", "account:B", 30)
    if err != nil {
        fmt.Println("Transfer failed:", err)
    } else {
        fmt.Println("Transfer successful")
    }
}
```

### 4.3 分布式锁 / Distributed Lock

```go
import (
    "github.com/go-redsync/redsync/v4"
    "github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

// 使用 redsync 实现分布式锁 / Distributed lock with redsync
func distributedLock(db *gorm.DB, rdb *redis.Client) {
    pool := goredis.NewPool(rdb)
    rs := redsync.New(pool)

    mutex := rs.NewMutex("resource-lock",
        redsync.WithExpiry(30*time.Second),
        redsync.WithTries(3),
    )

    // 获取锁
    if err := mutex.Lock(); err != nil {
        fmt.Println("Failed to acquire lock:", err)
        return
    }
    defer mutex.Unlock()

    // 执行需要互斥的操作 (如更新数据库)
    db.Model(&User{}).Where("id = ?", 1).Update("balance", gorm.Expr("balance - ?", 10))
    fmt.Println("Updated with lock")
}

// 手动实现分布式锁 / Manual implementation
type DistLock struct {
    rdb    *redis.Client
    key    string
    value  string
    expiry time.Duration
}

func NewDistLock(rdb *redis.Client, key string, expiry time.Duration) *DistLock {
    return &DistLock{
        rdb:    rdb,
        key:    key,
        value:  uuid.New().String(),
        expiry: expiry,
    }
}

func (l *DistLock) Lock(ctx context.Context) (bool, error) {
    return l.rdb.SetNX(ctx, l.key, l.value, l.expiry).Result()
}

func (l *DistLock) Unlock(ctx context.Context) error {
    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        end
        return 0
    `)
    _, err := script.Run(ctx, l.rdb, []string{l.key}, l.value).Result()
    return err
}
```

### 4.4 Lua 脚本 / Lua Scripts

```go
func luaScriptExample(ctx context.Context, rdb *redis.Client) {
    // 限流脚本
    rateLimitScript := redis.NewScript(`
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local window = tonumber(ARGV[2])
        
        local current = redis.call("INCR", key)
        if current == 1 then
            redis.call("EXPIRE", key, window)
        end
        
        if current > limit then
            return 0
        end
        return 1
    `)

    // 检查限流 (每分钟10次)
    for i := 0; i < 15; i++ {
        allowed, _ := rateLimitScript.Run(ctx, rdb, []string{"rate:user:1"}, 10, 60).Int64()
        if allowed == 1 {
            fmt.Printf("Request %d: Allowed\n", i+1)
        } else {
            fmt.Printf("Request %d: Rate limited\n", i+1)
        }
    }
}
```

## 5. 缓存模式 / Caching Patterns

### 5.1 Cache Aside 模式

```go
type CacheAside[T any] struct {
    db       *gorm.DB
    redis    *redis.Client
    prefix   string
    ttl      time.Duration
    loadFunc func(db *gorm.DB, key string) (T, error)
}

func (c *CacheAside[T]) Get(ctx context.Context, key string) (T, error) {
    var result T
    cacheKey := c.prefix + key

    // 1. 先查缓存
    data, err := c.redis.Get(ctx, cacheKey).Bytes()
    if err == nil {
        json.Unmarshal(data, &result)
        return result, nil
    }

    // 2. 缓存未命中，从数据库加载
    result, err = c.loadFunc(c.db, key)
    if err != nil {
        return result, err
    }

    // 3. 写入缓存
    data, _ = json.Marshal(result)
    c.redis.Set(ctx, cacheKey, data, c.ttl)

    return result, nil
}

func (c *CacheAside[T]) Delete(ctx context.Context, key string) error {
    return c.redis.Del(ctx, c.prefix+key).Err()
}
```

### 5.2 缓存穿透防护 / Cache Penetration Protection

```go
// 使用空值缓存防止缓存穿透
func (r *UserRepository) GetByIDWithNullCache(ctx context.Context, id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 检查缓存
    data, err := r.redis.Get(ctx, cacheKey).Bytes()
    if err == nil {
        if string(data) == "NULL" {
            return nil, gorm.ErrRecordNotFound
        }
        var user User
        json.Unmarshal(data, &user)
        return &user, nil
    }

    // 从数据库查询
    var user User
    if err := r.db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 缓存空值，防止穿透
            r.redis.Set(ctx, cacheKey, "NULL", 5*time.Minute)
        }
        return nil, err
    }

    // 写入缓存
    data, _ := json.Marshal(user)
    r.redis.Set(ctx, cacheKey, data, r.cacheTTL)

    return &user, nil
}
```

### 5.3 缓存雪崩防护 / Cache Avalanche Protection

```go
import "math/rand"

// 添加随机过期时间防止雪崩
func SetWithRandomExpiry(ctx context.Context, rdb *redis.Client, key string, value interface{}, baseTTL time.Duration) error {
    randomSeconds := time.Duration(rand.Intn(300)) * time.Second
    ttl := baseTTL + randomSeconds

    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return rdb.Set(ctx, key, data, ttl).Err()
}
```

## 6. 完整使用示例 / Complete Usage Example

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Name      string    `gorm:"size:100" json:"name"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    Age       int       `json:"age"`
    CreatedAt time.Time `json:"created_at"`
}

type UserService struct {
    db    *gorm.DB
    redis *redis.Client
}

func NewUserService(db *gorm.DB, rdb *redis.Client) *UserService {
    return &UserService{db: db, redis: rdb}
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)

    // 查缓存
    data, err := s.redis.Get(ctx, cacheKey).Bytes()
    if err == nil {
        var user User
        json.Unmarshal(data, &user)
        log.Println("Cache hit for user:", id)
        return &user, nil
    }

    // 查数据库
    var user User
    if err := s.db.First(&user, id).Error; err != nil {
        return nil, err
    }

    // 写缓存
    data, _ = json.Marshal(user)
    s.redis.Set(ctx, cacheKey, data, 10*time.Minute)
    log.Println("Cache miss, loaded from DB:", id)

    return &user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    if err := s.db.Create(user).Error; err != nil {
        return err
    }

    // 缓存新用户
    data, _ := json.Marshal(user)
    s.redis.Set(ctx, fmt.Sprintf("user:%d", user.ID), data, 10*time.Minute)

    return nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
    if err := s.db.Save(user).Error; err != nil {
        return err
    }

    // 删除旧缓存
    s.redis.Del(ctx, fmt.Sprintf("user:%d", user.ID))

    return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
    // 删除缓存
    s.redis.Del(ctx, fmt.Sprintf("user:%d", id))

    // 删除数据库记录
    return s.db.Delete(&User{}, id).Error
}

func main() {
    ctx := context.Background()

    // 连接 MySQL (GORM)
    dsn := "user:pass@tcp(localhost:3306)/dbname?parseTime=true"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect MySQL:", err)
    }
    sqlDB, _ := db.DB()
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)

    // 连接 Redis
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    defer rdb.Close()

    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatal("Failed to connect Redis:", err)
    }

    // 自动迁移
    db.AutoMigrate(&User{})

    // 使用服务
    userService := NewUserService(db, rdb)

    // 创建用户
    user := &User{Name: "Alice", Email: "alice@example.com", Age: 25}
    userService.CreateUser(ctx, user)
    fmt.Println("Created user:", user.ID)

    // 获取用户 (第一次从 DB)
    u, _ := userService.GetUser(ctx, user.ID)
    fmt.Printf("Got user: %+v\n", u)

    // 获取用户 (第二次从缓存)
    u, _ = userService.GetUser(ctx, user.ID)
    fmt.Printf("Got user (cached): %+v\n", u)

    // 更新用户
    user.Age = 26
    userService.UpdateUser(ctx, user)

    // 删除用户
    userService.DeleteUser(ctx, user.ID)
}
```
