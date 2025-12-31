# Redis 详解 / Redis Guide for Go

## 1. 连接与配置 / Connection & Configuration

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

func main() {
    ctx := context.Background()

    // ========================================
    // 单机连接 / Standalone Connection
    // ========================================
    rdb := redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        Password:     "",  // 无密码 / no password
        DB:           0,   // 默认 DB / default DB
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolSize:     10,
        MinIdleConns: 5,
        PoolTimeout:  4 * time.Second,
    })
    defer rdb.Close()

    // 测试连接 / Test connection
    pong, err := rdb.Ping(ctx).Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("Connected:", pong)

    // ========================================
    // 哨兵模式 / Sentinel Mode
    // ========================================
    rdbSentinel := redis.NewFailoverClient(&redis.FailoverOptions{
        MasterName:    "mymaster",
        SentinelAddrs: []string{"localhost:26379", "localhost:26380", "localhost:26381"},
        Password:      "",
        DB:            0,
    })
    defer rdbSentinel.Close()

    // ========================================
    // 集群模式 / Cluster Mode
    // ========================================
    rdbCluster := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs: []string{
            "localhost:7000",
            "localhost:7001",
            "localhost:7002",
        },
        Password:     "",
        PoolSize:     10,
        MinIdleConns: 5,
    })
    defer rdbCluster.Close()

    _ = rdbSentinel
    _ = rdbCluster
}
```

## 2. 基本数据类型操作 / Basic Data Type Operations

### 2.1 String 字符串

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
)

func stringOperations(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // SET / GET - 基本操作
    // ========================================
    
    // Set - 设置值 / Set value
    err := rdb.Set(ctx, "name", "Alice", 0).Err()
    if err != nil {
        panic(err)
    }

    // Set with expiration - 设置带过期时间的值
    err = rdb.Set(ctx, "token", "abc123", 10*time.Minute).Err()

    // SetNX - 仅当 key 不存在时设置 / Set only if not exists
    wasSet, err := rdb.SetNX(ctx, "lock", "1", 30*time.Second).Result()
    fmt.Println("SetNX result:", wasSet)

    // SetXX - 仅当 key 存在时设置 / Set only if exists
    wasSet, err = rdb.SetXX(ctx, "name", "Bob", 0).Result()
    fmt.Println("SetXX result:", wasSet)

    // Get - 获取值 / Get value
    val, err := rdb.Get(ctx, "name").Result()
    if err == redis.Nil {
        fmt.Println("Key does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("name:", val)
    }

    // GetSet - 获取旧值并设置新值 / Get old value and set new
    oldVal, err := rdb.GetSet(ctx, "name", "Charlie").Result()
    fmt.Println("Old value:", oldVal)

    // GetEx - 获取并设置过期时间 / Get and set expiration
    val, err = rdb.GetEx(ctx, "name", 1*time.Hour).Result()

    // GetDel - 获取并删除 / Get and delete
    val, err = rdb.GetDel(ctx, "token").Result()

    // ========================================
    // 批量操作 / Batch Operations
    // ========================================

    // MSet - 批量设置 / Batch set
    err = rdb.MSet(ctx, "key1", "val1", "key2", "val2", "key3", "val3").Err()

    // MGet - 批量获取 / Batch get
    vals, err := rdb.MGet(ctx, "key1", "key2", "key3", "nonexistent").Result()
    fmt.Println("MGet:", vals)  // [val1 val2 val3 <nil>]

    // ========================================
    // 数值操作 / Numeric Operations
    // ========================================

    rdb.Set(ctx, "counter", 0, 0)

    // Incr - 自增 1 / Increment by 1
    newVal, err := rdb.Incr(ctx, "counter").Result()
    fmt.Println("After Incr:", newVal)

    // IncrBy - 自增指定值 / Increment by specified value
    newVal, err = rdb.IncrBy(ctx, "counter", 10).Result()
    fmt.Println("After IncrBy(10):", newVal)

    // IncrByFloat - 浮点数自增 / Float increment
    floatVal, err := rdb.IncrByFloat(ctx, "counter", 0.5).Result()
    fmt.Println("After IncrByFloat(0.5):", floatVal)

    // Decr - 自减 1 / Decrement by 1
    newVal, err = rdb.Decr(ctx, "counter").Result()

    // DecrBy - 自减指定值 / Decrement by specified value
    newVal, err = rdb.DecrBy(ctx, "counter", 5).Result()

    // ========================================
    // 字符串操作 / String Operations
    // ========================================

    rdb.Set(ctx, "message", "Hello", 0)

    // Append - 追加 / Append
    newLen, err := rdb.Append(ctx, "message", " World").Result()
    fmt.Println("After Append, length:", newLen)

    // StrLen - 获取长度 / Get length
    length, err := rdb.StrLen(ctx, "message").Result()
    fmt.Println("StrLen:", length)

    // GetRange - 获取子串 / Get substring
    substr, err := rdb.GetRange(ctx, "message", 0, 4).Result()
    fmt.Println("GetRange(0,4):", substr)  // "Hello"

    // SetRange - 替换子串 / Replace substring
    newLen, err = rdb.SetRange(ctx, "message", 6, "Redis").Result()
    val, _ = rdb.Get(ctx, "message").Result()
    fmt.Println("After SetRange:", val)  // "Hello Redis"
}
```

### 2.2 Hash 哈希

```go
func hashOperations(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 基本操作 / Basic Operations
    // ========================================

    // HSet - 设置字段 / Set field
    err := rdb.HSet(ctx, "user:1", "name", "Alice", "age", 25, "city", "Beijing").Err()

    // HGet - 获取字段 / Get field
    name, err := rdb.HGet(ctx, "user:1", "name").Result()
    fmt.Println("Name:", name)

    // HGetAll - 获取所有字段 / Get all fields
    user, err := rdb.HGetAll(ctx, "user:1").Result()
    fmt.Println("User:", user)  // map[name:Alice age:25 city:Beijing]

    // HMGet - 批量获取 / Batch get
    vals, err := rdb.HMGet(ctx, "user:1", "name", "age").Result()
    fmt.Println("HMGet:", vals)

    // HExists - 检查字段是否存在 / Check if field exists
    exists, err := rdb.HExists(ctx, "user:1", "name").Result()
    fmt.Println("HExists:", exists)

    // HDel - 删除字段 / Delete field
    deleted, err := rdb.HDel(ctx, "user:1", "city").Result()
    fmt.Println("HDel count:", deleted)

    // HLen - 获取字段数量 / Get field count
    length, err := rdb.HLen(ctx, "user:1").Result()
    fmt.Println("HLen:", length)

    // HKeys - 获取所有字段名 / Get all field names
    keys, err := rdb.HKeys(ctx, "user:1").Result()
    fmt.Println("HKeys:", keys)

    // HVals - 获取所有字段值 / Get all field values
    values, err := rdb.HVals(ctx, "user:1").Result()
    fmt.Println("HVals:", values)

    // ========================================
    // 数值操作 / Numeric Operations
    // ========================================

    // HIncrBy - 整数自增 / Integer increment
    newAge, err := rdb.HIncrBy(ctx, "user:1", "age", 1).Result()
    fmt.Println("New age:", newAge)

    // HIncrByFloat - 浮点数自增 / Float increment
    rdb.HSet(ctx, "user:1", "score", 95.5)
    newScore, err := rdb.HIncrByFloat(ctx, "user:1", "score", 0.5).Result()
    fmt.Println("New score:", newScore)

    // ========================================
    // 条件设置 / Conditional Set
    // ========================================

    // HSetNX - 仅当字段不存在时设置 / Set only if field not exists
    wasSet, err := rdb.HSetNX(ctx, "user:1", "email", "alice@example.com").Result()
    fmt.Println("HSetNX result:", wasSet)

    // ========================================
    // 扫描 / Scan
    // ========================================

    // HScan - 增量迭代 / Incremental iteration
    iter := rdb.HScan(ctx, "user:1", 0, "*", 10).Iterator()
    for iter.Next(ctx) {
        fmt.Println("HScan:", iter.Val())
    }
}

// 使用 Hash 存储结构体 / Store struct in Hash
type User struct {
    ID    int64  `redis:"id"`
    Name  string `redis:"name"`
    Email string `redis:"email"`
    Age   int    `redis:"age"`
}

func structWithHash(ctx context.Context, rdb *redis.Client) {
    user := User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
        Age:   25,
    }

    // 使用 HSet 存储结构体 / Store struct with HSet
    err := rdb.HSet(ctx, "user:1", map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "age":   user.Age,
    }).Err()

    // 使用 HGetAll 读取到结构体 / Read to struct with HGetAll
    var loadedUser User
    err = rdb.HGetAll(ctx, "user:1").Scan(&loadedUser)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Loaded user: %+v\n", loadedUser)
}
```

### 2.3 List 列表

```go
func listOperations(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 基本操作 / Basic Operations
    // ========================================

    // LPush - 左侧插入 / Insert from left
    err := rdb.LPush(ctx, "queue", "a", "b", "c").Err()  // 结果: c, b, a

    // RPush - 右侧插入 / Insert from right
    err = rdb.RPush(ctx, "queue", "d", "e").Err()  // 结果: c, b, a, d, e

    // LPop - 左侧弹出 / Pop from left
    val, err := rdb.LPop(ctx, "queue").Result()
    fmt.Println("LPop:", val)  // "c"

    // RPop - 右侧弹出 / Pop from right
    val, err = rdb.RPop(ctx, "queue").Result()
    fmt.Println("RPop:", val)  // "e"

    // LLen - 获取长度 / Get length
    length, err := rdb.LLen(ctx, "queue").Result()
    fmt.Println("LLen:", length)

    // LRange - 获取范围 / Get range
    vals, err := rdb.LRange(ctx, "queue", 0, -1).Result()
    fmt.Println("LRange:", vals)

    // LIndex - 获取索引位置的元素 / Get element at index
    val, err = rdb.LIndex(ctx, "queue", 0).Result()
    fmt.Println("LIndex(0):", val)

    // LSet - 设置索引位置的元素 / Set element at index
    err = rdb.LSet(ctx, "queue", 0, "new_value").Err()

    // LInsert - 在元素前/后插入 / Insert before/after element
    err = rdb.LInsertBefore(ctx, "queue", "b", "inserted").Err()
    err = rdb.LInsertAfter(ctx, "queue", "b", "inserted_after").Err()

    // LRem - 移除元素 / Remove elements
    // count > 0: 从头到尾移除 count 个
    // count < 0: 从尾到头移除 |count| 个
    // count = 0: 移除所有
    removed, err := rdb.LRem(ctx, "queue", 0, "inserted").Result()
    fmt.Println("LRem count:", removed)

    // LTrim - 保留指定范围 / Keep specified range
    err = rdb.LTrim(ctx, "queue", 0, 2).Err()

    // ========================================
    // 阻塞操作 / Blocking Operations
    // ========================================

    // BLPop - 阻塞左侧弹出 / Blocking left pop
    result, err := rdb.BLPop(ctx, 5*time.Second, "queue").Result()
    if err == redis.Nil {
        fmt.Println("Timeout, no element")
    } else {
        fmt.Println("BLPop:", result)  // [queue value]
    }

    // BRPop - 阻塞右侧弹出 / Blocking right pop
    result, err = rdb.BRPop(ctx, 5*time.Second, "queue").Result()

    // BLMove - 阻塞移动元素 / Blocking move element
    val, err = rdb.BLMove(ctx, "queue1", "queue2", "LEFT", "RIGHT", 5*time.Second).Result()

    // ========================================
    // 消息队列模式 / Message Queue Pattern
    // ========================================

    // 生产者 / Producer
    go func() {
        for i := 0; i < 10; i++ {
            rdb.RPush(ctx, "tasks", fmt.Sprintf("task-%d", i))
            time.Sleep(100 * time.Millisecond)
        }
    }()

    // 消费者 / Consumer
    go func() {
        for {
            result, err := rdb.BLPop(ctx, 0, "tasks").Result()
            if err != nil {
                break
            }
            fmt.Println("Processing:", result[1])
        }
    }()
}
```

### 2.4 Set 集合

```go
func setOperations(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 基本操作 / Basic Operations
    // ========================================

    // SAdd - 添加成员 / Add members
    added, err := rdb.SAdd(ctx, "tags", "go", "redis", "docker", "kubernetes").Result()
    fmt.Println("SAdd count:", added)

    // SMembers - 获取所有成员 / Get all members
    members, err := rdb.SMembers(ctx, "tags").Result()
    fmt.Println("SMembers:", members)

    // SIsMember - 检查是否为成员 / Check if member
    isMember, err := rdb.SIsMember(ctx, "tags", "go").Result()
    fmt.Println("SIsMember 'go':", isMember)

    // SMIsMember - 批量检查 / Batch check
    results, err := rdb.SMIsMember(ctx, "tags", "go", "java", "redis").Result()
    fmt.Println("SMIsMember:", results)  // [true false true]

    // SCard - 获取成员数量 / Get member count
    count, err := rdb.SCard(ctx, "tags").Result()
    fmt.Println("SCard:", count)

    // SRem - 移除成员 / Remove members
    removed, err := rdb.SRem(ctx, "tags", "docker").Result()
    fmt.Println("SRem count:", removed)

    // SPop - 随机弹出成员 / Random pop member
    val, err := rdb.SPop(ctx, "tags").Result()
    fmt.Println("SPop:", val)

    // SPopN - 随机弹出 N 个成员 / Random pop N members
    vals, err := rdb.SPopN(ctx, "tags", 2).Result()
    fmt.Println("SPopN:", vals)

    // SRandMember - 随机获取成员(不移除) / Random get member (no remove)
    val, err = rdb.SRandMember(ctx, "tags").Result()

    // SRandMemberN - 随机获取 N 个成员 / Random get N members
    vals, err = rdb.SRandMemberN(ctx, "tags", 2).Result()

    // ========================================
    // 集合运算 / Set Operations
    // ========================================

    rdb.SAdd(ctx, "set1", "a", "b", "c", "d")
    rdb.SAdd(ctx, "set2", "c", "d", "e", "f")

    // SUnion - 并集 / Union
    union, err := rdb.SUnion(ctx, "set1", "set2").Result()
    fmt.Println("SUnion:", union)  // [a b c d e f]

    // SUnionStore - 并集并存储 / Union and store
    count, err = rdb.SUnionStore(ctx, "set_union", "set1", "set2").Result()

    // SInter - 交集 / Intersection
    inter, err := rdb.SInter(ctx, "set1", "set2").Result()
    fmt.Println("SInter:", inter)  // [c d]

    // SInterStore - 交集并存储 / Intersection and store
    count, err = rdb.SInterStore(ctx, "set_inter", "set1", "set2").Result()

    // SInterCard - 交集基数 / Intersection cardinality (Redis 7.0+)
    count, err = rdb.SInterCard(ctx, 0, "set1", "set2").Result()
    fmt.Println("SInterCard:", count)

    // SDiff - 差集 / Difference
    diff, err := rdb.SDiff(ctx, "set1", "set2").Result()
    fmt.Println("SDiff:", diff)  // [a b]

    // SDiffStore - 差集并存储 / Difference and store
    count, err = rdb.SDiffStore(ctx, "set_diff", "set1", "set2").Result()

    // SMove - 移动成员 / Move member
    moved, err := rdb.SMove(ctx, "set1", "set2", "a").Result()
    fmt.Println("SMove:", moved)

    // ========================================
    // 扫描 / Scan
    // ========================================

    iter := rdb.SScan(ctx, "tags", 0, "*", 10).Iterator()
    for iter.Next(ctx) {
        fmt.Println("SScan:", iter.Val())
    }
}
```

### 2.5 Sorted Set 有序集合

```go
func sortedSetOperations(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 基本操作 / Basic Operations
    // ========================================

    // ZAdd - 添加成员 / Add members
    added, err := rdb.ZAdd(ctx, "leaderboard", redis.Z{
        Score:  100,
        Member: "player1",
    }, redis.Z{
        Score:  85,
        Member: "player2",
    }, redis.Z{
        Score:  92,
        Member: "player3",
    }).Result()
    fmt.Println("ZAdd count:", added)

    // ZAddNX - 仅当成员不存在时添加 / Add only if not exists
    rdb.ZAddNX(ctx, "leaderboard", redis.Z{Score: 80, Member: "player4"})

    // ZAddXX - 仅当成员存在时更新 / Update only if exists
    rdb.ZAddXX(ctx, "leaderboard", redis.Z{Score: 105, Member: "player1"})

    // ZAddGT - 仅当新分数大于当前分数时更新 / Update only if new score > current
    rdb.ZAddGT(ctx, "leaderboard", redis.Z{Score: 110, Member: "player1"})

    // ZAddLT - 仅当新分数小于当前分数时更新 / Update only if new score < current
    rdb.ZAddLT(ctx, "leaderboard", redis.Z{Score: 80, Member: "player1"})

    // ZScore - 获取分数 / Get score
    score, err := rdb.ZScore(ctx, "leaderboard", "player1").Result()
    fmt.Println("ZScore player1:", score)

    // ZRank - 获取排名(升序) / Get rank (ascending)
    rank, err := rdb.ZRank(ctx, "leaderboard", "player1").Result()
    fmt.Println("ZRank player1:", rank)

    // ZRevRank - 获取排名(降序) / Get rank (descending)
    rank, err = rdb.ZRevRank(ctx, "leaderboard", "player1").Result()
    fmt.Println("ZRevRank player1:", rank)

    // ZCard - 获取成员数量 / Get member count
    count, err := rdb.ZCard(ctx, "leaderboard").Result()
    fmt.Println("ZCard:", count)

    // ZCount - 获取分数范围内的成员数量 / Count members in score range
    count, err = rdb.ZCount(ctx, "leaderboard", "80", "100").Result()
    fmt.Println("ZCount 80-100:", count)

    // ZIncrBy - 增加分数 / Increment score
    newScore, err := rdb.ZIncrBy(ctx, "leaderboard", 5, "player1").Result()
    fmt.Println("ZIncrBy player1:", newScore)

    // ========================================
    // 范围查询 / Range Queries
    // ========================================

    // ZRange - 按排名范围获取(升序) / Get by rank range (ascending)
    members, err := rdb.ZRange(ctx, "leaderboard", 0, -1).Result()
    fmt.Println("ZRange:", members)

    // ZRangeWithScores - 带分数 / With scores
    membersWithScores, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
    for _, z := range membersWithScores {
        fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
    }

    // ZRevRange - 按排名范围获取(降序) / Get by rank range (descending)
    members, err = rdb.ZRevRange(ctx, "leaderboard", 0, 2).Result()
    fmt.Println("Top 3:", members)

    // ZRangeByScore - 按分数范围获取 / Get by score range
    members, err = rdb.ZRangeByScore(ctx, "leaderboard", &redis.ZRangeBy{
        Min:    "80",
        Max:    "100",
        Offset: 0,
        Count:  10,
    }).Result()
    fmt.Println("ZRangeByScore 80-100:", members)

    // ZRevRangeByScore - 按分数范围获取(降序) / Get by score range (descending)
    members, err = rdb.ZRevRangeByScore(ctx, "leaderboard", &redis.ZRangeBy{
        Min: "-inf",
        Max: "+inf",
    }).Result()

    // ========================================
    // 删除操作 / Remove Operations
    // ========================================

    // ZRem - 移除成员 / Remove members
    removed, err := rdb.ZRem(ctx, "leaderboard", "player4").Result()

    // ZRemRangeByRank - 按排名范围移除 / Remove by rank range
    removed, err = rdb.ZRemRangeByRank(ctx, "leaderboard", 0, 0).Result()

    // ZRemRangeByScore - 按分数范围移除 / Remove by score range
    removed, err = rdb.ZRemRangeByScore(ctx, "leaderboard", "-inf", "60").Result()

    // ========================================
    // 集合运算 / Set Operations
    // ========================================

    rdb.ZAdd(ctx, "zset1", redis.Z{Score: 1, Member: "a"}, redis.Z{Score: 2, Member: "b"})
    rdb.ZAdd(ctx, "zset2", redis.Z{Score: 3, Member: "b"}, redis.Z{Score: 4, Member: "c"})

    // ZUnionStore - 并集 / Union
    count, err = rdb.ZUnionStore(ctx, "zset_union", &redis.ZStore{
        Keys:    []string{"zset1", "zset2"},
        Weights: []float64{1, 2},  // 权重
    }).Result()

    // ZInterStore - 交集 / Intersection
    count, err = rdb.ZInterStore(ctx, "zset_inter", &redis.ZStore{
        Keys: []string{"zset1", "zset2"},
    }).Result()
}
```

## 3. 高级功能 / Advanced Features

### 3.1 Pipeline 管道

```go
func pipelineExample(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 基本 Pipeline / Basic Pipeline
    // ========================================

    pipe := rdb.Pipeline()

    // 添加命令 / Add commands
    incr := pipe.Incr(ctx, "pipeline_counter")
    pipe.Expire(ctx, "pipeline_counter", time.Hour)
    get := pipe.Get(ctx, "pipeline_counter")

    // 执行 / Execute
    _, err := pipe.Exec(ctx)
    if err != nil {
        panic(err)
    }

    fmt.Println("Incr result:", incr.Val())
    fmt.Println("Get result:", get.Val())

    // ========================================
    // 使用 Pipelined 简化 / Simplified with Pipelined
    // ========================================

    var getValue *redis.StringCmd
    _, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
        pipe.Set(ctx, "key1", "value1", 0)
        pipe.Set(ctx, "key2", "value2", 0)
        getValue = pipe.Get(ctx, "key1")
        return nil
    })
    fmt.Println("Pipelined get:", getValue.Val())

    // ========================================
    // 批量操作 / Batch Operations
    // ========================================

    // 批量写入
    pipe = rdb.Pipeline()
    for i := 0; i < 1000; i++ {
        pipe.Set(ctx, fmt.Sprintf("batch_key_%d", i), i, 0)
    }
    _, err = pipe.Exec(ctx)

    // 批量读取
    pipe = rdb.Pipeline()
    gets := make([]*redis.StringCmd, 1000)
    for i := 0; i < 1000; i++ {
        gets[i] = pipe.Get(ctx, fmt.Sprintf("batch_key_%d", i))
    }
    _, err = pipe.Exec(ctx)

    for i, get := range gets {
        if i < 5 {  // 只打印前5个
            fmt.Printf("batch_key_%d = %s\n", i, get.Val())
        }
    }
}
```

### 3.2 Transaction 事务

```go
func transactionExample(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // TxPipeline - 事务管道 / Transaction Pipeline
    // ========================================

    pipe := rdb.TxPipeline()

    incr := pipe.Incr(ctx, "tx_counter")
    pipe.Expire(ctx, "tx_counter", time.Hour)

    _, err := pipe.Exec(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println("TxPipeline incr:", incr.Val())

    // ========================================
    // Watch - 乐观锁 / Optimistic Locking
    // ========================================

    // 转账示例 / Transfer example
    transfer := func(from, to string, amount int64) error {
        // 使用 WATCH 监视 key
        return rdb.Watch(ctx, func(tx *redis.Tx) error {
            // 获取当前余额
            fromBalance, err := tx.Get(ctx, from).Int64()
            if err != nil && err != redis.Nil {
                return err
            }

            if fromBalance < amount {
                return fmt.Errorf("insufficient balance")
            }

            toBalance, err := tx.Get(ctx, to).Int64()
            if err != nil && err != redis.Nil {
                return err
            }

            // 执行事务
            _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
                pipe.Set(ctx, from, fromBalance-amount, 0)
                pipe.Set(ctx, to, toBalance+amount, 0)
                return nil
            })
            return err
        }, from, to)  // 监视这两个 key
    }

    // 初始化余额
    rdb.Set(ctx, "account:A", 100, 0)
    rdb.Set(ctx, "account:B", 50, 0)

    // 转账
    err = transfer("account:A", "account:B", 30)
    if err != nil {
        fmt.Println("Transfer failed:", err)
    } else {
        fmt.Println("Transfer successful")
    }

    balanceA, _ := rdb.Get(ctx, "account:A").Result()
    balanceB, _ := rdb.Get(ctx, "account:B").Result()
    fmt.Printf("A: %s, B: %s\n", balanceA, balanceB)
}
```

### 3.3 Pub/Sub 发布订阅

```go
func pubsubExample(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 订阅 / Subscribe
    // ========================================

    pubsub := rdb.Subscribe(ctx, "channel1", "channel2")
    defer pubsub.Close()

    // 等待订阅确认 / Wait for subscription confirmation
    _, err := pubsub.Receive(ctx)
    if err != nil {
        panic(err)
    }

    // 方式1: 使用 Channel / Using Channel
    go func() {
        ch := pubsub.Channel()
        for msg := range ch {
            fmt.Printf("Received on %s: %s\n", msg.Channel, msg.Payload)
        }
    }()

    // 方式2: 使用 ReceiveMessage / Using ReceiveMessage
    go func() {
        for {
            msg, err := pubsub.ReceiveMessage(ctx)
            if err != nil {
                return
            }
            fmt.Printf("ReceiveMessage: %s: %s\n", msg.Channel, msg.Payload)
        }
    }()

    // ========================================
    // 发布 / Publish
    // ========================================

    time.Sleep(100 * time.Millisecond)  // 等待订阅就绪

    err = rdb.Publish(ctx, "channel1", "Hello Channel 1").Err()
    err = rdb.Publish(ctx, "channel2", "Hello Channel 2").Err()

    time.Sleep(100 * time.Millisecond)

    // ========================================
    // 模式订阅 / Pattern Subscribe
    // ========================================

    pubsubPattern := rdb.PSubscribe(ctx, "news.*")
    defer pubsubPattern.Close()

    go func() {
        ch := pubsubPattern.Channel()
        for msg := range ch {
            fmt.Printf("Pattern match %s on %s: %s\n", msg.Pattern, msg.Channel, msg.Payload)
        }
    }()

    time.Sleep(100 * time.Millisecond)

    rdb.Publish(ctx, "news.sports", "Sports news")
    rdb.Publish(ctx, "news.tech", "Tech news")

    time.Sleep(100 * time.Millisecond)
}
```

### 3.4 Lua 脚本 / Lua Scripts

```go
func luaScriptExample(ctx context.Context, rdb *redis.Client) {
    // ========================================
    // 基本脚本 / Basic Script
    // ========================================

    // 简单脚本
    script := redis.NewScript(`
        local key = KEYS[1]
        local increment = ARGV[1]
        local current = redis.call("GET", key) or "0"
        local new_value = tonumber(current) + tonumber(increment)
        redis.call("SET", key, new_value)
        return new_value
    `)

    result, err := script.Run(ctx, rdb, []string{"lua_counter"}, 10).Int64()
    if err != nil {
        panic(err)
    }
    fmt.Println("Lua script result:", result)

    // ========================================
    // 限流脚本 / Rate Limiting Script
    // ========================================

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

    // 检查限流
    for i := 0; i < 15; i++ {
        allowed, _ := rateLimitScript.Run(ctx, rdb, []string{"rate:user:1"}, 10, 60).Int64()
        if allowed == 1 {
            fmt.Printf("Request %d: Allowed\n", i+1)
        } else {
            fmt.Printf("Request %d: Rate limited\n", i+1)
        }
    }

    // ========================================
    // 分布式锁脚本 / Distributed Lock Script
    // ========================================

    // 获取锁
    acquireLockScript := redis.NewScript(`
        if redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2]) then
            return 1
        end
        return 0
    `)

    // 释放锁
    releaseLockScript := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        end
        return 0
    `)

    lockKey := "lock:resource"
    lockValue := "unique-id-123"
    lockTTL := 30000  // 30 seconds

    acquired, _ := acquireLockScript.Run(ctx, rdb, []string{lockKey}, lockValue, lockTTL).Int64()
    if acquired == 1 {
        fmt.Println("Lock acquired")
        
        // 执行业务逻辑...
        
        released, _ := releaseLockScript.Run(ctx, rdb, []string{lockKey}, lockValue).Int64()
        if released == 1 {
            fmt.Println("Lock released")
        }
    }
}
```

### 3.5 分布式锁 / Distributed Lock

```go
import (
    "github.com/go-redsync/redsync/v4"
    "github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

func distributedLockExample(ctx context.Context, rdb *redis.Client) {
    // 使用 redsync 库 / Using redsync library
    pool := goredis.NewPool(rdb)
    rs := redsync.New(pool)

    // 创建互斥锁 / Create mutex
    mutex := rs.NewMutex("my-lock",
        redsync.WithExpiry(30*time.Second),
        redsync.WithTries(3),
        redsync.WithRetryDelay(100*time.Millisecond),
    )

    // 获取锁 / Acquire lock
    if err := mutex.Lock(); err != nil {
        fmt.Println("Failed to acquire lock:", err)
        return
    }
    fmt.Println("Lock acquired")

    // 执行业务逻辑 / Execute business logic
    time.Sleep(2 * time.Second)

    // 释放锁 / Release lock
    if ok, err := mutex.Unlock(); !ok || err != nil {
        fmt.Println("Failed to release lock:", err)
        return
    }
    fmt.Println("Lock released")

    // ========================================
    // 手动实现分布式锁 / Manual Implementation
    // ========================================

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

    func (l *DistLock) Extend(ctx context.Context, duration time.Duration) (bool, error) {
        script := redis.NewScript(`
            if redis.call("GET", KEYS[1]) == ARGV[1] then
                return redis.call("PEXPIRE", KEYS[1], ARGV[2])
            end
            return 0
        `)
        result, err := script.Run(ctx, l.rdb, []string{l.key}, l.value, int(duration.Milliseconds())).Int64()
        return result == 1, err
    }
}
```

## 4. 缓存模式 / Caching Patterns

```go
package cache

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    "github.com/redis/go-redis/v9"
)

// ========================================
// Cache Aside 模式 / Cache Aside Pattern
// ========================================

type CacheAside[T any] struct {
    rdb    *redis.Client
    prefix string
    ttl    time.Duration
    loader func(ctx context.Context, key string) (T, error)
}

func NewCacheAside[T any](rdb *redis.Client, prefix string, ttl time.Duration, loader func(ctx context.Context, key string) (T, error)) *CacheAside[T] {
    return &CacheAside[T]{
        rdb:    rdb,
        prefix: prefix,
        ttl:    ttl,
        loader: loader,
    }
}

func (c *CacheAside[T]) Get(ctx context.Context, key string) (T, error) {
    var result T
    cacheKey := c.prefix + key

    // 1. 先查缓存 / Check cache first
    data, err := c.rdb.Get(ctx, cacheKey).Bytes()
    if err == nil {
        if err := json.Unmarshal(data, &result); err == nil {
            return result, nil
        }
    }

    // 2. 缓存未命中，从数据源加载 / Cache miss, load from source
    if err == redis.Nil || err != nil {
        result, err = c.loader(ctx, key)
        if err != nil {
            return result, err
        }

        // 3. 写入缓存 / Write to cache
        data, _ := json.Marshal(result)
        c.rdb.Set(ctx, cacheKey, data, c.ttl)
    }

    return result, nil
}

func (c *CacheAside[T]) Delete(ctx context.Context, key string) error {
    return c.rdb.Del(ctx, c.prefix+key).Err()
}

// ========================================
// 缓存穿透防护 / Cache Penetration Protection
// ========================================

type CacheWithBloom[T any] struct {
    rdb        *redis.Client
    bloomKey   string
    cacheAside *CacheAside[T]
}

func (c *CacheWithBloom[T]) Get(ctx context.Context, key string) (T, error) {
    var zero T

    // 检查布隆过滤器 / Check bloom filter
    exists, _ := c.rdb.BFExists(ctx, c.bloomKey, key).Result()
    if !exists {
        return zero, errors.New("key not exists")
    }

    return c.cacheAside.Get(ctx, key)
}

// ========================================
// 缓存雪崩防护 / Cache Avalanche Protection
// ========================================

func SetWithRandomExpiry(ctx context.Context, rdb *redis.Client, key string, value interface{}, baseTTL time.Duration) error {
    // 添加随机时间避免同时过期
    // Add random time to avoid simultaneous expiration
    randomSeconds := time.Duration(rand.Intn(300)) * time.Second
    ttl := baseTTL + randomSeconds

    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return rdb.Set(ctx, key, data, ttl).Err()
}

// ========================================
// 缓存击穿防护 (使用 singleflight) / Cache Breakdown Protection
// ========================================

import "golang.org/x/sync/singleflight"

type CacheWithSingleFlight[T any] struct {
    rdb    *redis.Client
    prefix string
    ttl    time.Duration
    loader func(ctx context.Context, key string) (T, error)
    sf     singleflight.Group
}

func (c *CacheWithSingleFlight[T]) Get(ctx context.Context, key string) (T, error) {
    var result T
    cacheKey := c.prefix + key

    // 1. 先查缓存
    data, err := c.rdb.Get(ctx, cacheKey).Bytes()
    if err == nil {
        json.Unmarshal(data, &result)
        return result, nil
    }

    // 2. 使用 singleflight 防止缓存击穿
    v, err, _ := c.sf.Do(key, func() (interface{}, error) {
        // 再次检查缓存 (double check)
        data, err := c.rdb.Get(ctx, cacheKey).Bytes()
        if err == nil {
            var r T
            json.Unmarshal(data, &r)
            return r, nil
        }

        // 从数据源加载
        result, err := c.loader(ctx, key)
        if err != nil {
            return nil, err
        }

        // 写入缓存
        data, _ := json.Marshal(result)
        c.rdb.Set(ctx, cacheKey, data, c.ttl)

        return result, nil
    })

    if err != nil {
        return result, err
    }

    return v.(T), nil
}
```

## 5. 使用示例 / Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

func main() {
    ctx := context.Background()

    // 创建客户端 / Create client
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    defer rdb.Close()

    // 测试连接 / Test connection
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    fmt.Println("Connected to Redis!")

    // 基本操作示例 / Basic operations example
    rdb.Set(ctx, "greeting", "Hello, Redis!", 10*time.Minute)
    
    val, err := rdb.Get(ctx, "greeting").Result()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("greeting:", val)

    // 哈希示例 / Hash example
    rdb.HSet(ctx, "user:1", map[string]interface{}{
        "name":  "Alice",
        "email": "alice@example.com",
        "age":   25,
    })

    user, _ := rdb.HGetAll(ctx, "user:1").Result()
    fmt.Println("User:", user)

    // 列表示例 / List example
    rdb.RPush(ctx, "queue", "task1", "task2", "task3")
    
    tasks, _ := rdb.LRange(ctx, "queue", 0, -1).Result()
    fmt.Println("Tasks:", tasks)

    // 集合示例 / Set example
    rdb.SAdd(ctx, "tags", "go", "redis", "docker")
    
    tags, _ := rdb.SMembers(ctx, "tags").Result()
    fmt.Println("Tags:", tags)

    // 有序集合示例 / Sorted set example
    rdb.ZAdd(ctx, "scores", 
        redis.Z{Score: 100, Member: "player1"},
        redis.Z{Score: 85, Member: "player2"},
        redis.Z{Score: 92, Member: "player3"},
    )
    
    top3, _ := rdb.ZRevRangeWithScores(ctx, "scores", 0, 2).Result()
    fmt.Println("Top 3:")
    for i, z := range top3 {
        fmt.Printf("  %d. %s: %.0f\n", i+1, z.Member, z.Score)
    }
}
```