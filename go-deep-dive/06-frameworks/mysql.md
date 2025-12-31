# MySQL 操作 (使用 GORM) / MySQL Operations (with GORM)

## 1. 连接与配置 / Connection & Configuration

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

// 安装 / Installation:
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql

func main() {
    // DSN 格式 / DSN format:
    // [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...]
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

    // 打开连接 / Open connection
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),  // 日志级别
    })
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }

    // 获取底层 *sql.DB 以配置连接池 / Get underlying *sql.DB to configure pool
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get DB:", err)
    }

    // 配置连接池 / Configure connection pool
    sqlDB.SetMaxOpenConns(25)                  // 最大打开连接数
    sqlDB.SetMaxIdleConns(5)                   // 最大空闲连接数
    sqlDB.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生存时间
    sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // 空闲连接最大时间

    // 验证连接 / Verify connection
    if err := sqlDB.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }

    fmt.Println("Database connected!")
}

// 高级连接配置 / Advanced connection configuration
func advancedConnection() (*gorm.DB, error) {
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

    return gorm.Open(mysql.New(mysql.Config{
        DSN:                       dsn,
        DefaultStringSize:         256,   // 默认字符串长度
        DisableDatetimePrecision:  true,  // 禁用 datetime 精度
        DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式
        DontSupportRenameColumn:   true,  // 用 `change` 重命名列
        SkipInitializeWithVersion: false, // 根据版本自动配置
    }), &gorm.Config{
        Logger:                 logger.Default.LogMode(logger.Info),
        SkipDefaultTransaction: true,  // 跳过默认事务以提高性能
        PrepareStmt:            true,  // 缓存预编译语句
    })
}
```

## 2. 模型定义与自动迁移 / Model Definition & Auto Migration

```go
package main

import (
    "database/sql"
    "time"

    "gorm.io/gorm"
)

// 用户模型 / User model
type User struct {
    ID        uint           `gorm:"primarykey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"size:255;uniqueIndex"`
    Age       int            `gorm:"default:0"`
    Active    bool           `gorm:"default:true"`
    Profile   Profile        `gorm:"foreignKey:UserID"`  // has one
    Posts     []Post         `gorm:"foreignKey:UserID"`  // has many
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`  // 软删除
}

type Profile struct {
    ID     uint
    UserID uint
    Bio    string `gorm:"type:text"`
    Avatar string `gorm:"size:255"`
}

type Post struct {
    ID        uint
    UserID    uint
    Title     string `gorm:"size:200;not null"`
    Content   string `gorm:"type:text"`
    Published bool   `gorm:"default:false"`
    CreatedAt time.Time
}

// 自动迁移 / Auto migration
func autoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(&User{}, &Profile{}, &Post{})
}
```

## 3. 基本 CRUD 操作 / Basic CRUD Operations

```go
package main

import (
    "fmt"
    "gorm.io/gorm"
)

// Create 创建 / Create
func createUser(db *gorm.DB, name, email string, age int) (*User, error) {
    user := &User{
        Name:  name,
        Email: email,
        Age:   age,
    }
    result := db.Create(user)
    return user, result.Error
}

// 批量创建 / Batch create
func batchCreateUsers(db *gorm.DB, users []User) error {
    return db.CreateInBatches(users, 100).Error
}

// Read 查询单条 / Query single row
func getUserByID(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.First(&user, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

// 条件查询 / Conditional query
func getUserByEmail(db *gorm.DB, email string) (*User, error) {
    var user User
    result := db.Where("email = ?", email).First(&user)
    return &user, result.Error
}

// Read 查询多条 / Query multiple rows
func getUsers(db *gorm.DB, minAge int) ([]User, error) {
    var users []User
    result := db.Where("age >= ?", minAge).Find(&users)
    return users, result.Error
}

// 带预加载的查询 / Query with preloading
func getUserWithPosts(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.Preload("Posts").Preload("Profile").First(&user, id)
    return &user, result.Error
}

// Update 更新 / Update
func updateUser(db *gorm.DB, id uint, name string, age int) error {
    return db.Model(&User{}).Where("id = ?", id).Updates(map[string]interface{}{
        "name": name,
        "age":  age,
    }).Error
}

// 更新单个字段 / Update single field
func updateUserAge(db *gorm.DB, id uint, age int) error {
    return db.Model(&User{}).Where("id = ?", id).Update("age", age).Error
}

// 使用结构体更新 / Update with struct
func updateUserByStruct(db *gorm.DB, id uint, updates User) error {
    return db.Model(&User{ID: id}).Updates(updates).Error
}

// Delete 删除 / Delete (软删除)
func deleteUser(db *gorm.DB, id uint) error {
    return db.Delete(&User{}, id).Error
}

// 永久删除 / Hard delete
func hardDeleteUser(db *gorm.DB, id uint) error {
    return db.Unscoped().Delete(&User{}, id).Error
}

// 查询已删除记录 / Query soft deleted records
func getDeletedUsers(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
    return users, result.Error
}
```

## 4. 事务处理 / Transaction Handling

```go
package main

import (
    "fmt"
    "gorm.io/gorm"
)

// 基本事务 / Basic transaction
func transferMoney(db *gorm.DB, fromID, toID uint, amount float64) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 扣款 / Deduct
        if err := tx.Model(&Account{}).Where("id = ?", fromID).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return err
        }

        // 入账 / Credit
        if err := tx.Model(&Account{}).Where("id = ?", toID).
            Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
            return err
        }

        return nil
    })
}

type Account struct {
    ID      uint
    UserID  uint
    Balance float64
}

// 手动事务控制 / Manual transaction control
func manualTransaction(db *gorm.DB) error {
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    if err := tx.Create(&User{Name: "Alice"}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Create(&User{Name: "Bob"}).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

// 嵌套事务 / Nested transaction
func nestedTransaction(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        tx.Create(&User{Name: "User1"})

        // 嵌套事务 (使用保存点)
        tx.Transaction(func(tx2 *gorm.DB) error {
            tx2.Create(&User{Name: "User2"})
            return nil  // 提交保存点
        })

        return nil  // 提交外层事务
    })
}
```

## 5. 高级查询 / Advanced Queries

```go
package main

import (
    "gorm.io/gorm"
)

// 动态构建查询 / Dynamic query building
func searchUsers(db *gorm.DB, name string, minAge, maxAge int, active *bool) ([]User, error) {
    query := db.Model(&User{})

    if name != "" {
        query = query.Where("name LIKE ?", "%"+name+"%")
    }

    if minAge > 0 {
        query = query.Where("age >= ?", minAge)
    }

    if maxAge > 0 {
        query = query.Where("age <= ?", maxAge)
    }

    if active != nil {
        query = query.Where("active = ?", *active)
    }

    var users []User
    result := query.Find(&users)
    return users, result.Error
}

// IN 查询 / IN query
func getUsersByIDs(db *gorm.DB, ids []uint) ([]User, error) {
    var users []User
    result := db.Where("id IN ?", ids).Find(&users)
    return users, result.Error
}

// OR 查询 / OR query
func searchByNameOrEmail(db *gorm.DB, keyword string) ([]User, error) {
    var users []User
    result := db.Where("name LIKE ?", "%"+keyword+"%").
        Or("email LIKE ?", "%"+keyword+"%").
        Find(&users)
    return users, result.Error
}

// 分页查询 / Pagination query
func getUsersPaginated(db *gorm.DB, page, pageSize int) ([]User, int64, error) {
    var users []User
    var total int64

    // 获取总数 / Get total count
    db.Model(&User{}).Count(&total)

    // 查询当前页 / Query current page
    offset := (page - 1) * pageSize
    result := db.Offset(offset).Limit(pageSize).Find(&users)

    return users, total, result.Error
}

// 排序 / Ordering
func getUsersSorted(db *gorm.DB, field string, desc bool) ([]User, error) {
    var users []User
    order := field
    if desc {
        order += " DESC"
    }
    result := db.Order(order).Find(&users)
    return users, result.Error
}

// 选择特定字段 / Select specific fields
func getUserNames(db *gorm.DB) ([]string, error) {
    var names []string
    result := db.Model(&User{}).Pluck("name", &names)
    return names, result.Error
}

// 聚合查询 / Aggregate queries
type UserStats struct {
    Count  int64
    AvgAge float64
    MaxAge int
    MinAge int
}

func getUserStats(db *gorm.DB) (*UserStats, error) {
    var stats UserStats
    result := db.Model(&User{}).Where("active = ?", true).Select(
        "COUNT(*) as count",
        "AVG(age) as avg_age",
        "MAX(age) as max_age",
        "MIN(age) as min_age",
    ).Scan(&stats)
    return &stats, result.Error
}

// 分组查询 / Group by query
type AgeGroup struct {
    Age   int
    Count int64
}

func getUserCountByAge(db *gorm.DB) ([]AgeGroup, error) {
    var groups []AgeGroup
    result := db.Model(&User{}).
        Select("age, COUNT(*) as count").
        Group("age").
        Having("COUNT(*) > ?", 1).
        Scan(&groups)
    return groups, result.Error
}

// 子查询 / Subquery
func getUsersAboveAvgAge(db *gorm.DB) ([]User, error) {
    var users []User
    subQuery := db.Model(&User{}).Select("AVG(age)")
    result := db.Where("age > (?)", subQuery).Find(&users)
    return users, result.Error
}

// JOIN 查询 / JOIN query
func getUsersWithPosts(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Joins("JOIN posts ON posts.user_id = users.id").
        Where("posts.published = ?", true).
        Distinct().
        Find(&users)
    return users, result.Error
}

// 原生 SQL / Raw SQL
func rawQuery(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Raw("SELECT * FROM users WHERE age > ? AND active = ?", 18, true).Scan(&users)
    return users, result.Error
}

// 原生 SQL 执行 / Raw SQL execution
func rawExec(db *gorm.DB) error {
    return db.Exec("UPDATE users SET age = age + 1 WHERE id = ?", 1).Error
}
```

## 6. Hooks (钩子) / Hooks

```go
package main

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "time"

    "gorm.io/gorm"
)

type UserWithHooks struct {
    ID           uint
    Name         string
    Email        string
    Password     string
    PasswordHash string `gorm:"size:64"`
    LastLoginAt  *time.Time
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// BeforeCreate 创建前钩子 / Before create hook
func (u *UserWithHooks) BeforeCreate(tx *gorm.DB) error {
    // 验证 / Validation
    if u.Name == "" {
        return errors.New("name cannot be empty")
    }

    // 密码哈希 / Password hashing
    if u.Password != "" {
        hash := sha256.Sum256([]byte(u.Password))
        u.PasswordHash = hex.EncodeToString(hash[:])
        u.Password = ""  // 清除明文密码
    }

    return nil
}

// AfterCreate 创建后钩子 / After create hook
func (u *UserWithHooks) AfterCreate(tx *gorm.DB) error {
    // 可以执行一些创建后的操作，如发送欢迎邮件
    // Can perform post-create operations like sending welcome email
    return nil
}

// BeforeUpdate 更新前钩子 / Before update hook
func (u *UserWithHooks) BeforeUpdate(tx *gorm.DB) error {
    // 如果密码被更新，重新哈希
    if u.Password != "" {
        hash := sha256.Sum256([]byte(u.Password))
        u.PasswordHash = hex.EncodeToString(hash[:])
        u.Password = ""
    }
    return nil
}

// AfterFind 查询后钩子 / After find hook
func (u *UserWithHooks) AfterFind(tx *gorm.DB) error {
    // 可以执行一些查询后的处理
    return nil
}

// BeforeDelete 删除前钩子 / Before delete hook
func (u *UserWithHooks) BeforeDelete(tx *gorm.DB) error {
    // 可以执行删除前的检查
    return nil
}
```

## 7. 关联操作 / Association Operations

```go
package main

import (
    "gorm.io/gorm"
)

// 预加载 / Preloading
func getUserWithAssociations(db *gorm.DB, id uint) (*User, error) {
    var user User

    // 预加载所有关联 / Preload all associations
    result := db.Preload("Profile").Preload("Posts").First(&user, id)
    return &user, result.Error
}

// 条件预加载 / Conditional preloading
func getUserWithPublishedPosts(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.Preload("Posts", "published = ?", true).First(&user, id)
    return &user, result.Error
}

// 嵌套预加载 / Nested preloading
func getUserWithNestedAssociations(db *gorm.DB, id uint) (*User, error) {
    var user User
    result := db.Preload("Posts.Comments").First(&user, id)
    return &user, result.Error
}

// 添加关联 / Add association
func addPostToUser(db *gorm.DB, userID uint, post *Post) error {
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        return err
    }
    return db.Model(&user).Association("Posts").Append(post)
}

// 移除关联 / Remove association
func removePostFromUser(db *gorm.DB, userID uint, postID uint) error {
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        return err
    }
    var post Post
    if err := db.First(&post, postID).Error; err != nil {
        return err
    }
    return db.Model(&user).Association("Posts").Delete(&post)
}

// 替换关联 / Replace associations
func replaceUserPosts(db *gorm.DB, userID uint, posts []Post) error {
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        return err
    }
    return db.Model(&user).Association("Posts").Replace(posts)
}

// 清除关联 / Clear associations
func clearUserPosts(db *gorm.DB, userID uint) error {
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        return err
    }
    return db.Model(&user).Association("Posts").Clear()
}

// 计数关联 / Count associations
func countUserPosts(db *gorm.DB, userID uint) (int64, error) {
    var user User
    if err := db.First(&user, userID).Error; err != nil {
        return 0, err
    }
    return db.Model(&user).Association("Posts").Count(), nil
}
```

## 8. 性能优化与最佳实践 / Performance & Best Practices

```go
package main

import (
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

// 1. 跳过默认事务以提高性能 / Skip default transaction
func setupWithSkipTransaction() (*gorm.DB, error) {
    return gorm.Open(mysql.Open(dsn), &gorm.Config{
        SkipDefaultTransaction: true,
    })
}

// 2. 使用预编译语句 / Use prepared statements
func setupWithPrepareStmt() (*gorm.DB, error) {
    return gorm.Open(mysql.Open(dsn), &gorm.Config{
        PrepareStmt: true,
    })
}

// 3. 批量插入 / Batch insert
func batchInsert(db *gorm.DB, users []User) error {
    return db.CreateInBatches(users, 100).Error
}

// 4. 只查询需要的字段 / Select only needed fields
func getMinimalUsers(db *gorm.DB) ([]User, error) {
    var users []User
    result := db.Select("id", "name", "email").Find(&users)
    return users, result.Error
}

// 5. 使用索引提示 / Use index hints
func queryWithIndex(db *gorm.DB, email string) (*User, error) {
    var user User
    result := db.Clauses(clause.Index{Name: "idx_users_email"}).
        Where("email = ?", email).First(&user)
    return &user, result.Error
}

// 6. 使用 Find 而非 First 进行批量查询 / Use Find instead of First for batch queries
func getMultipleUsers(db *gorm.DB, ids []uint) ([]User, error) {
    var users []User
    result := db.Find(&users, ids)  // 比循环使用 First 更高效
    return users, result.Error
}

// 7. 使用 Map 更新避免零值问题 / Use Map for updates to avoid zero value issues
func updateWithMap(db *gorm.DB, id uint, updates map[string]interface{}) error {
    return db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}

// 8. 使用 Upsert / Use Upsert
func upsertUser(db *gorm.DB, user *User) error {
    return db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "email"}},
        DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
    }).Create(user).Error
}

// 9. 使用连接池监控 / Monitor connection pool
func monitorPool(db *gorm.DB) {
    sqlDB, _ := db.DB()
    stats := sqlDB.Stats()

    fmt.Printf("MaxOpenConnections: %d\n", stats.MaxOpenConnections)
    fmt.Printf("OpenConnections: %d\n", stats.OpenConnections)
    fmt.Printf("InUse: %d\n", stats.InUse)
    fmt.Printf("Idle: %d\n", stats.Idle)
    fmt.Printf("WaitCount: %d\n", stats.WaitCount)
    fmt.Printf("WaitDuration: %v\n", stats.WaitDuration)
}

// 10. 使用 DryRun 调试查询 / Use DryRun to debug queries
func debugQuery(db *gorm.DB) {
    stmt := db.Session(&gorm.Session{DryRun: true}).
        Where("age > ?", 18).
        Find(&User{}).Statement

    fmt.Println("SQL:", stmt.SQL.String())
    fmt.Println("Vars:", stmt.Vars)
}
```

## 9. 日志与调试 / Logging & Debugging

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// 自定义日志 / Custom logger
func setupWithCustomLogger() (*gorm.DB, error) {
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             200 * time.Millisecond,  // 慢查询阈值
            LogLevel:                  logger.Info,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )

    return gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
}

// 临时开启调试 / Temporarily enable debug
func debugSession(db *gorm.DB) *gorm.DB {
    return db.Debug()
}

// 记录慢查询 / Log slow queries
type SlowQueryLogger struct {
    logger.Interface
    threshold time.Duration
}

func (l *SlowQueryLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
    elapsed := time.Since(begin)
    sql, rows := fc()

    if elapsed > l.threshold {
        log.Printf("[SLOW QUERY] %s | %v | %d rows | %s", elapsed, err, rows, sql)
    }

    l.Interface.Trace(ctx, begin, fc, err)
}
```
