# Go 数据库与框架详解

## 1. GORM - Go ORM 框架

### 1.1 GORM 基础

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// 模型定义
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:200"`
    Age       int            `gorm:"default:0"`
    Active    bool           `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除
}

// 自定义表名
func (User) TableName() string {
    return "users"
}

// 钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
    log.Println("Before creating user:", u.Name)
    return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
    log.Println("After creating user:", u.ID)
    return nil
}

func main() {
    // 连接数据库
    dsn := "user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }

    // 获取底层 *sql.DB
    sqlDB, _ := db.DB()
    
    // 连接池配置
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // 自动迁移
    db.AutoMigrate(&User{})

    // CRUD 操作
    crudExample(db)
}

func crudExample(db *gorm.DB) {
    // Create - 创建
    user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
    result := db.Create(&user)
    fmt.Printf("Created user ID: %d, affected rows: %d\n", user.ID, result.RowsAffected)

    // 批量创建
    users := []User{
        {Name: "Bob", Email: "bob@example.com", Age: 30},
        {Name: "Charlie", Email: "charlie@example.com", Age: 35},
    }
    db.Create(&users)

    // Read - 查询
    var foundUser User
    db.First(&foundUser, 1) // 通过主键查找
    fmt.Printf("Found user: %+v\n", foundUser)

    db.First(&foundUser, "email = ?", "alice@example.com") // 条件查找

    // 查询多条
    var allUsers []User
    db.Find(&allUsers)
    fmt.Printf("All users count: %d\n", len(allUsers))

    // 条件查询
    var activeUsers []User
    db.Where("active = ? AND age > ?", true, 20).Find(&activeUsers)

    // Update - 更新
    db.Model(&user).Update("Age", 26)
    db.Model(&user).Updates(User{Name: "Alice Updated", Age: 27})
    db.Model(&user).Updates(map[string]interface{}{"name": "Alice", "age": 28})

    // Delete - 删除
    db.Delete(&user, 1) // 软删除

    // 永久删除
    db.Unscoped().Delete(&user, 1)
}
```

### 1.2 GORM 高级查询

```go
package main

import (
    "fmt"
    "time"

    "gorm.io/gorm"
)

type Product struct {
    ID        uint
    Name      string
    Price     float64
    Stock     int
    Category  string
    CreatedAt time.Time
}

func advancedQueries(db *gorm.DB) {
    // Select 指定字段
    var products []Product
    db.Select("name", "price").Find(&products)

    // Order 排序
    db.Order("price desc").Find(&products)
    db.Order("price desc, name").Find(&products)

    // Limit & Offset 分页
    db.Limit(10).Offset(0).Find(&products)

    // Group & Having
    type Result struct {
        Category string
        Total    int
    }
    var results []Result
    db.Model(&Product{}).
        Select("category, sum(stock) as total").
        Group("category").
        Having("total > ?", 100).
        Scan(&results)

    // Distinct
    var categories []string
    db.Model(&Product{}).Distinct().Pluck("category", &categories)

    // 子查询
    subQuery := db.Model(&Product{}).Select("AVG(price)")
    db.Where("price > (?)", subQuery).Find(&products)

    // 原生 SQL
    db.Raw("SELECT * FROM products WHERE price > ?", 100).Scan(&products)
    db.Exec("UPDATE products SET stock = stock - 1 WHERE id = ?", 1)

    // 链式条件
    query := db.Model(&Product{})
    if true { // 条件判断
        query = query.Where("price > ?", 50)
    }
    if true {
        query = query.Where("stock > ?", 0)
    }
    query.Find(&products)

    // Scopes 复用查询条件
    db.Scopes(PriceGreaterThan(100), InStock).Find(&products)

    // Joins
    type UserOrder struct {
        UserName  string
        OrderID   uint
        Amount    float64
    }
    var userOrders []UserOrder
    db.Table("users").
        Select("users.name as user_name, orders.id as order_id, orders.amount").
        Joins("left join orders on users.id = orders.user_id").
        Scan(&userOrders)

    // Preload 预加载
    type Order struct {
        ID        uint
        UserID    uint
        User      User
        Items     []OrderItem
    }
    type OrderItem struct {
        ID       uint
        OrderID  uint
        Product  string
        Quantity int
    }
    
    var orders []Order
    db.Preload("User").Preload("Items").Find(&orders)

    // 带条件的预加载
    db.Preload("Items", "quantity > ?", 2).Find(&orders)

    // 嵌套预加载
    db.Preload("Items.Product").Find(&orders)
}

// Scope 函数
func PriceGreaterThan(price float64) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("price > ?", price)
    }
}

func InStock(db *gorm.DB) *gorm.DB {
    return db.Where("stock > 0")
}

// 事务
func transactionExample(db *gorm.DB) error {
    // 方式1: 自动事务
    err := db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&Product{Name: "Product A", Price: 100}).Error; err != nil {
            return err
        }

        if err := tx.Create(&Product{Name: "Product B", Price: 200}).Error; err != nil {
            return err
        }

        return nil
    })

    // 方式2: 手动事务
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Create(&Product{Name: "Product C"}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Create(&Product{Name: "Product D"}).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error

    return err
}

// 关联关系
type Author struct {
    ID    uint
    Name  string
    Books []Book // 一对多
}

type Book struct {
    ID       uint
    Title    string
    AuthorID uint
    Author   Author // 属于
    Tags     []Tag  `gorm:"many2many:book_tags"` // 多对多
}

type Tag struct {
    ID    uint
    Name  string
    Books []Book `gorm:"many2many:book_tags"`
}

func associationExample(db *gorm.DB) {
    // 创建关联
    author := Author{
        Name: "John",
        Books: []Book{
            {Title: "Book 1"},
            {Title: "Book 2"},
        },
    }
    db.Create(&author)

    // 添加关联
    db.Model(&author).Association("Books").Append(&Book{Title: "Book 3"})

    // 查询关联
    var books []Book
    db.Model(&author).Association("Books").Find(&books)

    // 删除关联
    db.Model(&author).Association("Books").Delete(&books[0])

    // 清空关联
    db.Model(&author).Association("Books").Clear()

    // 替换关联
    db.Model(&author).Association("Books").Replace(&Book{Title: "New Book"})

    // 关联数量
    count := db.Model(&author).Association("Books").Count()
    fmt.Println("Books count:", count)
}
```

## 2. MySQL 与 GORM

### 2.1 MySQL 连接配置

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

type Config struct {
    Host         string
    Port         int
    User         string
    Password     string
    DBName       string
    MaxIdleConns int
    MaxOpenConns int
    MaxLifetime  time.Duration
}

func NewMySQLDB(cfg Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

    db, err := gorm.Open(mysql.New(mysql.Config{
        DSN:                       dsn,
        DefaultStringSize:         256,
        DisableDatetimePrecision:  true,
        DontSupportRenameIndex:    true,
        DontSupportRenameColumn:   true,
        SkipInitializeWithVersion: false,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NamingStrategy: schema.NamingStrategy{
            TablePrefix:   "t_",
            SingularTable: true,
        },
        DisableForeignKeyConstraintWhenMigrating: true,
    })

    if err != nil {
        return nil, err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }

    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

    return db, nil
}

func main() {
    cfg := Config{
        Host:         "127.0.0.1",
        Port:         3306,
        User:         "root",
        Password:     "password",
        DBName:       "testdb",
        MaxIdleConns: 10,
        MaxOpenConns: 100,
        MaxLifetime:  time.Hour,
    }

    db, err := NewMySQLDB(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // 测试连接
    sqlDB, _ := db.DB()
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("Database ping failed:", err)
    }

    log.Println("Database connected successfully")
}
```

### 2.2 MySQL 特定功能

```go
package main

import (
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type Article struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"size:200;not null;index"`
    Content   string    `gorm:"type:text"`
    Status    int       `gorm:"default:0;index"`
    ViewCount int       `gorm:"default:0"`
    CreatedAt time.Time `gorm:"index"`
    UpdatedAt time.Time
}

func mysqlFeatures(db *gorm.DB) {
    // Upsert - 插入或更新
    article := Article{ID: 1, Title: "Test", Content: "Content"}
    db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "id"}},
        DoUpdates: clause.AssignmentColumns([]string{"title", "content"}),
    }).Create(&article)

    // 批量插入优化
    var articles []Article
    for i := 0; i < 1000; i++ {
        articles = append(articles, Article{Title: "Article"})
    }
    db.CreateInBatches(&articles, 100) // 每批100条

    // 乐观锁
    type Product struct {
        ID      uint
        Name    string
        Version int `gorm:"default:0"`
    }
    
    var product Product
    db.First(&product, 1)
    
    db.Model(&product).
        Where("version = ?", product.Version).
        Updates(map[string]interface{}{
            "name":    "Updated",
            "version": product.Version + 1,
        })

    // 悲观锁
    tx := db.Begin()
    tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, 1)
    // ... 操作
    tx.Commit()

    // 全文索引
    db.Raw("SELECT * FROM articles WHERE MATCH(title, content) AGAINST(? IN BOOLEAN MODE)", "keyword").
        Scan(&articles)

    // 分区表查询
    db.Table("articles PARTITION (p2024)").Find(&articles)
}

// 数据库连接池监控
func monitorPool(db *gorm.DB) {
    sqlDB, _ := db.DB()
    
    stats := sqlDB.Stats()
    
    println("MaxOpenConnections:", stats.MaxOpenConnections)
    println("OpenConnections:", stats.OpenConnections)
    println("InUse:", stats.InUse)
    println("Idle:", stats.Idle)
    println("WaitCount:", stats.WaitCount)
    println("WaitDuration:", stats.WaitDuration)
}
```

## 3. Redis 集成

### 3.1 go-redis 基础

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
    // 单机连接
    rdb := redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,
        PoolSize:     10,
        MinIdleConns: 5,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })

    // 测试连接
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatal("Redis ping failed:", err)
    }
    fmt.Println("Ping:", pong)

    // 基础操作
    basicOperations(rdb)

    // 高级操作
    advancedOperations(rdb)

    // 关闭连接
    rdb.Close()
}

func basicOperations(rdb *redis.Client) {
    // String
    err := rdb.Set(ctx, "name", "Alice", 0).Err()
    if err != nil {
        log.Println("Set error:", err)
    }

    val, err := rdb.Get(ctx, "name").Result()
    if err == redis.Nil {
        fmt.Println("key does not exist")
    } else if err != nil {
        log.Println("Get error:", err)
    } else {
        fmt.Println("name:", val)
    }

    // 设置过期时间
    rdb.Set(ctx, "temp", "value", 10*time.Second)
    rdb.SetEx(ctx, "temp2", "value", 10*time.Second)

    // INCR/DECR
    rdb.Set(ctx, "counter", 0, 0)
    rdb.Incr(ctx, "counter")
    rdb.IncrBy(ctx, "counter", 5)
    rdb.Decr(ctx, "counter")

    // Hash
    rdb.HSet(ctx, "user:1", "name", "Bob", "age", 30)
    rdb.HGet(ctx, "user:1", "name")
    rdb.HGetAll(ctx, "user:1")
    rdb.HIncrBy(ctx, "user:1", "age", 1)

    // List
    rdb.LPush(ctx, "queue", "item1", "item2", "item3")
    rdb.RPush(ctx, "queue", "item4")
    rdb.LRange(ctx, "queue", 0, -1)
    rdb.LPop(ctx, "queue")
    rdb.RPop(ctx, "queue")
    rdb.BLPop(ctx, 5*time.Second, "queue") // 阻塞

    // Set
    rdb.SAdd(ctx, "tags", "go", "redis", "mysql")
    rdb.SMembers(ctx, "tags")
    rdb.SIsMember(ctx, "tags", "go")
    rdb.SRem(ctx, "tags", "mysql")

    // Sorted Set
    rdb.ZAdd(ctx, "leaderboard",
        redis.Z{Score: 100, Member: "Alice"},
        redis.Z{Score: 200, Member: "Bob"},
        redis.Z{Score: 150, Member: "Charlie"},
    )
    rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1)
    rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2) // Top 3
    rdb.ZIncrBy(ctx, "leaderboard", 10, "Alice")
}

func advancedOperations(rdb *redis.Client) {
    // Pipeline
    pipe := rdb.Pipeline()
    pipe.Set(ctx, "key1", "value1", 0)
    pipe.Set(ctx, "key2", "value2", 0)
    pipe.Get(ctx, "key1")
    cmds, _ := pipe.Exec(ctx)
    for _, cmd := range cmds {
        fmt.Println(cmd)
    }

    // 事务
    tx := rdb.TxPipeline()
    tx.Set(ctx, "tx_key1", "value1", 0)
    tx.Set(ctx, "tx_key2", "value2", 0)
    tx.Exec(ctx)

    // Watch (乐观锁)
    err := rdb.Watch(ctx, func(tx *redis.Tx) error {
        val, err := tx.Get(ctx, "counter").Int()
        if err != nil && err != redis.Nil {
            return err
        }

        _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
            pipe.Set(ctx, "counter", val+1, 0)
            return nil
        })
        return err
    }, "counter")
    if err == redis.TxFailedErr {
        fmt.Println("Transaction failed due to concurrent modification")
    }

    // Lua 脚本
    script := redis.NewScript(`
        local key = KEYS[1]
        local change = ARGV[1]
        local value = redis.call("GET", key)
        if not value then
            value = 0
        end
        value = value + change
        redis.call("SET", key, value)
        return value
    `)
    result, _ := script.Run(ctx, rdb, []string{"lua_counter"}, 10).Int()
    fmt.Println("Lua result:", result)

    // Pub/Sub
    go func() {
        pubsub := rdb.Subscribe(ctx, "channel1")
        defer pubsub.Close()

        for msg := range pubsub.Channel() {
            fmt.Println("Received:", msg.Channel, msg.Payload)
        }
    }()

    time.Sleep(100 * time.Millisecond)
    rdb.Publish(ctx, "channel1", "Hello, Redis!")
}

// 缓存封装
type Cache struct {
    client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
    return &Cache{client: client}
}

func (c *Cache) Get(key string, dest interface{}) error {
    val, err := c.client.Get(ctx, key).Result()
    if err != nil {
        return err
    }
    return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *Cache) Delete(key string) error {
    return c.client.Del(ctx, key).Err()
}
```

### 3.2 Redis 分布式锁

```go
package main

import (
    "context"
    "errors"
    "time"

    "github.com/redis/go-redis/v9"
)

type DistributedLock struct {
    client     *redis.Client
    key        string
    value      string
    expiration time.Duration
}

func NewDistributedLock(client *redis.Client, key string, expiration time.Duration) *DistributedLock {
    return &DistributedLock{
        client:     client,
        key:        "lock:" + key,
        value:      generateUUID(), // 唯一标识
        expiration: expiration,
    }
}

func (l *DistributedLock) Lock(ctx context.Context) (bool, error) {
    return l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
}

func (l *DistributedLock) TryLock(ctx context.Context, timeout time.Duration) (bool, error) {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        ok, err := l.Lock(ctx)
        if err != nil {
            return false, err
        }
        if ok {
            return true, nil
        }
        time.Sleep(50 * time.Millisecond)
    }
    return false, nil
}

func (l *DistributedLock) Unlock(ctx context.Context) error {
    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `)
    
    result, err := script.Run(ctx, l.client, []string{l.key}, l.value).Int()
    if err != nil {
        return err
    }
    if result == 0 {
        return errors.New("lock not held")
    }
    return nil
}

func (l *DistributedLock) Extend(ctx context.Context, extension time.Duration) (bool, error) {
    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("PEXPIRE", KEYS[1], ARGV[2])
        else
            return 0
        end
    `)
    
    result, err := script.Run(ctx, l.client, []string{l.key}, l.value, extension.Milliseconds()).Int()
    if err != nil {
        return false, err
    }
    return result == 1, nil
}

func generateUUID() string {
    // 简化实现，实际应使用 uuid 库
    return "unique-id-" + time.Now().String()
}

// 使用示例
func lockExample(rdb *redis.Client) {
    ctx := context.Background()
    lock := NewDistributedLock(rdb, "resource-1", 10*time.Second)

    // 获取锁
    ok, err := lock.TryLock(ctx, 5*time.Second)
    if err != nil {
        return
    }
    if !ok {
        return // 获取失败
    }

    // 确保释放锁
    defer lock.Unlock(ctx)

    // 执行业务逻辑
}
```

## 4. Thrift

### 4.1 Thrift 定义文件

```thrift
// user.thrift
namespace go user

struct User {
    1: required i64 id
    2: required string name
    3: optional string email
    4: optional i32 age
}

struct CreateUserRequest {
    1: required string name
    2: optional string email
    3: optional i32 age
}

struct CreateUserResponse {
    1: required User user
}

struct GetUserRequest {
    1: required i64 id
}

struct GetUserResponse {
    1: optional User user
}

struct ListUsersRequest {
    1: optional i32 page = 1
    2: optional i32 pageSize = 10
}

struct ListUsersResponse {
    1: required list<User> users
    2: required i32 total
}

exception UserNotFoundException {
    1: string message
}

service UserService {
    CreateUserResponse CreateUser(1: CreateUserRequest req)
    GetUserResponse GetUser(1: GetUserRequest req) throws (1: UserNotFoundException e)
    ListUsersResponse ListUsers(1: ListUsersRequest req)
}
```

### 4.2 Thrift 服务端

```go
package main

import (
    "context"
    "log"
    "sync"

    "github.com/apache/thrift/lib/go/thrift"
    "your-project/gen-go/user"
)

type UserServiceImpl struct {
    users  map[int64]*user.User
    nextID int64
    mu     sync.RWMutex
}

func NewUserServiceImpl() *UserServiceImpl {
    return &UserServiceImpl{
        users:  make(map[int64]*user.User),
        nextID: 1,
    }
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    newUser := &user.User{
        ID:    s.nextID,
        Name:  req.Name,
        Email: req.Email,
        Age:   req.Age,
    }
    s.users[s.nextID] = newUser
    s.nextID++

    return &user.CreateUserResponse{User: newUser}, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    u, ok := s.users[req.ID]
    if !ok {
        return nil, &user.UserNotFoundException{Message: "User not found"}
    }

    return &user.GetUserResponse{User: u}, nil
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var users []*user.User
    for _, u := range s.users {
        users = append(users, u)
    }

    return &user.ListUsersResponse{
        Users: users,
        Total: int32(len(users)),
    }, nil
}

func main() {
    handler := NewUserServiceImpl()
    processor := user.NewUserServiceProcessor(handler)

    transport, err := thrift.NewTServerSocket(":9090")
    if err != nil {
        log.Fatal(err)
    }

    transportFactory := thrift.NewTBufferedTransportFactory(8192)
    protocolFactory := thrift.NewTBinaryProtocolFactoryConf(nil)

    server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

    log.Println("Thrift server starting on :9090")
    if err := server.Serve(); err != nil {
        log.Fatal(err)
    }
}
```

## 5. Kitex - 字节跳动 RPC 框架

### 5.1 Kitex 服务端

```go
package main

import (
    "context"
    "log"

    "github.com/cloudwego/kitex/pkg/rpcinfo"
    "github.com/cloudwego/kitex/server"
    user "your-project/kitex_gen/user/userservice"
)

type UserServiceImpl struct{}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
    // 业务逻辑
    return &user.CreateUserResponse{
        User: &user.User{
            ID:   1,
            Name: req.Name,
        },
    }, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
    return &user.GetUserResponse{
        User: &user.User{
            ID:   req.ID,
            Name: "Test User",
        },
    }, nil
}

func main() {
    svr := user.NewServer(
        new(UserServiceImpl),
        server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
        server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
            ServiceName: "user.service",
        }),
    )

    if err := svr.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### 5.2 Kitex 客户端

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/cloudwego/kitex/client"
    "github.com/cloudwego/kitex/pkg/retry"
    user "your-project/kitex_gen/user/userservice"
)

func main() {
    // 创建客户端
    c, err := user.NewClient(
        "user.service",
        client.WithHostPorts("127.0.0.1:8888"),
        client.WithRPCTimeout(3*time.Second),
        client.WithConnectTimeout(time.Second),
        client.WithFailureRetry(retry.NewFailurePolicy()),
    )
    if err != nil {
        log.Fatal(err)
    }

    // 调用服务
    resp, err := c.CreateUser(context.Background(), &user.CreateUserRequest{
        Name: "Alice",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Created user: %+v", resp.User)
}
```

### 5.3 Kitex 中间件

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/cloudwego/kitex/pkg/endpoint"
)

// 日志中间件
func LoggingMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) error {
        start := time.Now()
        err := next(ctx, req, resp)
        log.Printf("Request took %v, error: %v", time.Since(start), err)
        return err
    }
}

// 恢复中间件
func RecoveryMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
    return func(ctx context.Context, req, resp interface{}) (err error) {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Panic recovered: %v", r)
                err = fmt.Errorf("internal error")
            }
        }()
        return next(ctx, req, resp)
    }
}

// 限流中间件
func RateLimitMiddleware(limiter *rate.Limiter) endpoint.Middleware {
    return func(next endpoint.Endpoint) endpoint.Endpoint {
        return func(ctx context.Context, req, resp interface{}) error {
            if !limiter.Allow() {
                return errors.New("rate limit exceeded")
            }
            return next(ctx, req, resp)
        }
    }
}
```

## 6. Hertz - 字节跳动 HTTP 框架

### 6.1 Hertz 基础

```go
package main

import (
    "context"
    "net/http"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
    h := server.Default(
        server.WithHostPorts(":8080"),
        server.WithReadTimeout(10*time.Second),
        server.WithWriteTimeout(10*time.Second),
    )

    // 基础路由
    h.GET("/", func(c context.Context, ctx *app.RequestContext) {
        ctx.JSON(consts.StatusOK, utils.H{
            "message": "Hello, Hertz!",
        })
    })

    // 路径参数
    h.GET("/users/:id", func(c context.Context, ctx *app.RequestContext) {
        id := ctx.Param("id")
        ctx.JSON(consts.StatusOK, utils.H{
            "id": id,
        })
    })

    // 查询参数
    h.GET("/search", func(c context.Context, ctx *app.RequestContext) {
        query := ctx.Query("q")
        page := ctx.DefaultQuery("page", "1")
        ctx.JSON(consts.StatusOK, utils.H{
            "query": query,
            "page":  page,
        })
    })

    // POST 请求
    h.POST("/users", func(c context.Context, ctx *app.RequestContext) {
        var user struct {
            Name  string `json:"name"`
            Email string `json:"email"`
        }

        if err := ctx.BindJSON(&user); err != nil {
            ctx.JSON(consts.StatusBadRequest, utils.H{"error": err.Error()})
            return
        }

        ctx.JSON(consts.StatusCreated, utils.H{
            "user": user,
        })
    })

    // 路由组
    v1 := h.Group("/api/v1")
    {
        v1.GET("/users", listUsers)
        v1.POST("/users", createUser)
        v1.PUT("/users/:id", updateUser)
        v1.DELETE("/users/:id", deleteUser)
    }

    h.Spin()
}

func listUsers(c context.Context, ctx *app.RequestContext) {
    ctx.JSON(consts.StatusOK, utils.H{"users": []string{}})
}

func createUser(c context.Context, ctx *app.RequestContext) {
    ctx.JSON(consts.StatusCreated, utils.H{"status": "created"})
}

func updateUser(c context.Context, ctx *app.RequestContext) {
    ctx.JSON(consts.StatusOK, utils.H{"status": "updated"})
}

func deleteUser(c context.Context, ctx *app.RequestContext) {
    ctx.JSON(consts.StatusOK, utils.H{"status": "deleted"})
}
```

### 6.2 Hertz 中间件

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
    "github.com/cloudwego/hertz/pkg/common/utils"
    "github.com/cloudwego/hertz/pkg/protocol/consts"
)

// 日志中间件
func LoggerMiddleware() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        start := time.Now()
        path := string(ctx.Request.URI().Path())
        method := string(ctx.Request.Method())

        ctx.Next(c)

        latency := time.Since(start)
        status := ctx.Response.StatusCode()

        log.Printf("[%s] %s %d %v", method, path, status, latency)
    }
}

// 认证中间件
func AuthMiddleware() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        token := string(ctx.GetHeader("Authorization"))
        if token == "" {
            ctx.AbortWithStatusJSON(consts.StatusUnauthorized, utils.H{
                "error": "unauthorized",
            })
            return
        }

        // 验证 token
        userID := validateToken(token)
        if userID == "" {
            ctx.AbortWithStatusJSON(consts.StatusUnauthorized, utils.H{
                "error": "invalid token",
            })
            return
        }

        ctx.Set("userID", userID)
        ctx.Next(c)
    }
}

func validateToken(token string) string {
    // 实际验证逻辑
    return "user123"
}

// CORS 中间件
func CORSMiddleware() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        ctx.Header("Access-Control-Allow-Origin", "*")
        ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if string(ctx.Method()) == "OPTIONS" {
            ctx.AbortWithStatus(consts.StatusNoContent)
            return
        }

        ctx.Next(c)
    }
}

// 恢复中间件
func RecoveryMiddleware() app.HandlerFunc {
    return func(c context.Context, ctx *app.RequestContext) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic: %v", err)
                ctx.AbortWithStatusJSON(consts.StatusInternalServerError, utils.H{
                    "error": "internal server error",
                })
            }
        }()

        ctx.Next(c)
    }
}

func main() {
    h := server.Default()

    // 全局中间件
    h.Use(RecoveryMiddleware())
    h.Use(LoggerMiddleware())
    h.Use(CORSMiddleware())

    // 公开路由
    h.GET("/health", func(c context.Context, ctx *app.RequestContext) {
        ctx.JSON(consts.StatusOK, utils.H{"status": "ok"})
    })

    // 需要认证的路由
    api := h.Group("/api")
    api.Use(AuthMiddleware())
    {
        api.GET("/profile", func(c context.Context, ctx *app.RequestContext) {
            userID := ctx.GetString("userID")
            ctx.JSON(consts.StatusOK, utils.H{"userID": userID})
        })
    }

    h.Spin()
}
```

## 7. Docker 与 Docker Compose

### 7.1 Dockerfile

```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

# 安装依赖
RUN apk add --no-cache git ca-certificates tzdata

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/server

# 运行阶段
FROM alpine:3.18

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN adduser -D -g '' appuser

# 工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
ENTRYPOINT ["./server"]
```

### 7.2 Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
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
      - DB_USER=app
      - DB_PASSWORD=secret
      - DB_NAME=appdb
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
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M

  mysql:
    image: mysql:8.0
    container_name: mysql
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=appdb
      - MYSQL_USER=app
      - MYSQL_PASSWORD=secret
    volumes:
      - mysql-data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped
    command: redis-server --appendonly yes

  nginx:
    image: nginx:alpine
    container_name: nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./certs:/etc/nginx/certs
    depends_on:
      - app
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

### 7.3 开发环境 Compose

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
      - "2345:2345"  # Delve 调试端口
    environment:
      - APP_ENV=development
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=root
      - DB_PASSWORD=root
      - DB_NAME=appdb
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - .:/app
      - go-mod:/go/pkg/mod
    depends_on:
      - mysql
      - redis
    networks:
      - dev-network
    command: air -c .air.toml

  mysql:
    image: mysql:8.0
    container_name: mysql-dev
    ports:
      - "3306:3306"
    environment:
      - MYSQL_ROOT_PASSWORD=root
      - MYSQL_DATABASE=appdb
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

  adminer:
    image: adminer
    container_name: adminer
    ports:
      - "8081:8080"
    networks:
      - dev-network

networks:
  dev-network:
    driver: bridge

volumes:
  mysql-dev-data:
  go-mod:
```

### 7.4 Dockerfile.dev

```dockerfile
# Dockerfile.dev
FROM golang:1.21

# 安装开发工具
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 热重载配置会在 volume 挂载时使用
EXPOSE 8080 2345

CMD ["air", "-c", ".air.toml"]
```

### 7.5 Air 配置 (.air.toml)

```toml
# .air.toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -gcflags='all=-N -l' -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

### 7.6 常用 Docker 命令

```bash
# 构建镜像
docker build -t myapp:latest .

# 运行容器
docker run -d -p 8080:8080 --name myapp myapp:latest

# 查看日志
docker logs -f myapp

# 进入容器
docker exec -it myapp sh

# Docker Compose 命令
docker-compose up -d
docker-compose down
docker-compose logs -f app
docker-compose ps
docker-compose restart app

# 清理
docker system prune -a
docker volume prune
```
