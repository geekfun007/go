# GORM - Go ORM 框架 / Go ORM Framework

## 1. 安装与连接 / Installation & Connection

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "gorm.io/driver/mysql"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// 安装 / Installation:
// go get -u gorm.io/gorm
// go get -u gorm.io/driver/mysql
// go get -u gorm.io/driver/sqlite

func main() {
    // SQLite 连接 / SQLite connection
    dbSqlite, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to SQLite:", err)
    }
    _ = dbSqlite
    
    // MySQL 连接 / MySQL connection
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal("Failed to connect to MySQL:", err)
    }
    
    // 配置连接池 / Configure connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
    sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
    
    fmt.Println("Database connected!")
}
```

## 2. 模型定义 / Model Definition

```go
package main

import (
    "time"
    
    "gorm.io/gorm"
)

// 基础模型 / Basic model
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Name      string         `gorm:"size:100;not null"`
    Email     string         `gorm:"uniqueIndex;size:255"`
    Age       int            `gorm:"default:18"`
    Birthday  *time.Time
    Active    bool           `gorm:"default:true"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`  // 软删除
}

// 自定义表名 / Custom table name
func (User) TableName() string {
    return "users"
}

// 带关联的模型 / Model with associations

// 一对一 / One-to-one
type Profile struct {
    ID     uint
    UserID uint
    User   User `gorm:"foreignKey:UserID"`
    Bio    string
    Avatar string
}

// 一对多 / One-to-many
type Post struct {
    ID        uint
    Title     string
    Content   string
    UserID    uint
    User      User       `gorm:"foreignKey:UserID"`
    Comments  []Comment  `gorm:"foreignKey:PostID"`
    Tags      []Tag      `gorm:"many2many:post_tags;"`  // 多对多
    CreatedAt time.Time
}

type Comment struct {
    ID        uint
    Content   string
    PostID    uint
    UserID    uint
    User      User `gorm:"foreignKey:UserID"`
    CreatedAt time.Time
}

// 多对多 / Many-to-many
type Tag struct {
    ID    uint
    Name  string `gorm:"uniqueIndex"`
    Posts []Post `gorm:"many2many:post_tags;"`
}

// 自定义字段类型 / Custom field types
type JSON map[string]interface{}

type Config struct {
    ID       uint
    Settings JSON `gorm:"type:json"`
}

// 嵌入结构体 / Embedded struct
type BaseModel struct {
    ID        uint `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Article struct {
    BaseModel  // 嵌入基础模型
    Title   string
    Content string
}

// 自动迁移 / Auto migration
func autoMigrate(db *gorm.DB) {
    db.AutoMigrate(&User{}, &Profile{}, &Post{}, &Comment{}, &Tag{})
}
```

## 3. CRUD 操作 / CRUD Operations

```go
package main

import (
    "fmt"
    "time"
    
    "gorm.io/gorm"
)

// Create 创建 / Create
func createExamples(db *gorm.DB) {
    // 创建单条记录 / Create single record
    user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
    result := db.Create(&user)
    
    fmt.Println("ID:", user.ID)                // 自动填充的 ID
    fmt.Println("Rows affected:", result.RowsAffected)
    fmt.Println("Error:", result.Error)
    
    // 批量创建 / Batch create
    users := []User{
        {Name: "Bob", Email: "bob@example.com"},
        {Name: "Charlie", Email: "charlie@example.com"},
    }
    db.Create(&users)
    
    // 指定字段创建 / Create with selected fields
    db.Select("Name", "Email").Create(&User{
        Name:  "Diana",
        Email: "diana@example.com",
        Age:   30,  // 会被忽略
    })
    
    // 排除字段创建 / Create omitting fields
    db.Omit("Age").Create(&User{
        Name:  "Eve",
        Email: "eve@example.com",
        Age:   30,  // 会被忽略
    })
    
    // 使用 Map 创建 / Create with map
    db.Model(&User{}).Create(map[string]interface{}{
        "Name":  "Frank",
        "Email": "frank@example.com",
    })
}

// Read 查询 / Read
func readExamples(db *gorm.DB) {
    var user User
    var users []User
    
    // 查询第一条 / First record
    db.First(&user)  // 按主键排序
    fmt.Println("First:", user)
    
    // 查询最后一条 / Last record
    db.Last(&user)
    
    // 按主键查询 / Find by primary key
    db.First(&user, 1)                    // 主键为 1
    db.First(&user, "id = ?", 1)          // 等价写法
    db.Find(&users, []int{1, 2, 3})       // 多个主键
    
    // 条件查询 / Query with conditions
    db.Where("name = ?", "Alice").First(&user)
    db.Where("age > ?", 20).Find(&users)
    db.Where("name LIKE ?", "%li%").Find(&users)
    db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)
    db.Where("name IN ?", []string{"Alice", "Bob"}).Find(&users)
    
    // 结构体条件 / Struct conditions
    db.Where(&User{Name: "Alice", Age: 25}).First(&user)
    
    // Map 条件 / Map conditions
    db.Where(map[string]interface{}{"name": "Alice", "age": 25}).First(&user)
    
    // Not 条件 / Not conditions
    db.Not("name = ?", "Alice").Find(&users)
    db.Not(map[string]interface{}{"name": "Alice"}).Find(&users)
    
    // Or 条件 / Or conditions
    db.Where("name = ?", "Alice").Or("name = ?", "Bob").Find(&users)
    
    // 选择字段 / Select fields
    db.Select("name", "email").Find(&users)
    db.Select("name", "age as user_age").Find(&users)
    
    // 排序 / Order
    db.Order("age desc").Find(&users)
    db.Order("age desc, name asc").Find(&users)
    
    // 分页 / Pagination
    db.Limit(10).Offset(0).Find(&users)   // 第1页，每页10条
    db.Limit(10).Offset(10).Find(&users)  // 第2页
    
    // 分组 / Group
    type Result struct {
        Age   int
        Count int
    }
    var results []Result
    db.Model(&User{}).Select("age, count(*) as count").Group("age").Find(&results)
    
    // Distinct
    db.Distinct("name", "age").Find(&users)
    
    // Count
    var count int64
    db.Model(&User{}).Where("age > ?", 20).Count(&count)
    
    // Pluck - 查询单列 / Query single column
    var names []string
    db.Model(&User{}).Pluck("name", &names)
    
    // 扫描到 struct / Scan to struct
    type UserInfo struct {
        Name  string
        Email string
    }
    var userInfo UserInfo
    db.Model(&User{}).Select("name", "email").First(&userInfo)
}

// Update 更新 / Update
func updateExamples(db *gorm.DB) {
    var user User
    db.First(&user, 1)
    
    // 更新单个字段 / Update single field
    db.Model(&user).Update("name", "NewName")
    
    // 更新多个字段 / Update multiple fields
    db.Model(&user).Updates(User{Name: "NewName", Age: 30})
    db.Model(&user).Updates(map[string]interface{}{"name": "NewName", "age": 30})
    
    // 选择字段更新 / Select fields to update
    db.Model(&user).Select("name").Updates(User{Name: "NewName", Age: 30})  // 只更新 name
    
    // 排除字段更新 / Omit fields
    db.Model(&user).Omit("name").Updates(User{Name: "NewName", Age: 30})  // 不更新 name
    
    // 批量更新 / Batch update
    db.Model(&User{}).Where("age < ?", 18).Updates(map[string]interface{}{"active": false})
    
    // 使用 SQL 表达式 / Use SQL expression
    db.Model(&user).Update("age", gorm.Expr("age + ?", 1))
    
    // UpdateColumn 跳过 Hooks / Skip Hooks
    db.Model(&user).UpdateColumn("name", "NewName")
}

// Delete 删除 / Delete
func deleteExamples(db *gorm.DB) {
    var user User
    
    // 删除单条 / Delete single
    db.Delete(&user, 1)  // DELETE FROM users WHERE id = 1
    
    // 条件删除 / Delete with conditions
    db.Where("name = ?", "Alice").Delete(&User{})
    db.Delete(&User{}, "age < ?", 18)
    
    // 批量删除 / Batch delete
    db.Delete(&User{}, []int{1, 2, 3})
    
    // 软删除 (需要 DeletedAt 字段) / Soft delete
    db.Delete(&user)  // 只设置 DeletedAt
    
    // 查询包含软删除的记录 / Find soft deleted records
    db.Unscoped().Where("id = ?", 1).Find(&user)
    
    // 永久删除 / Permanent delete
    db.Unscoped().Delete(&user)
}
```

## 4. 关联操作 / Association Operations

```go
package main

import "gorm.io/gorm"

func associationExamples(db *gorm.DB) {
    // 创建带关联的记录 / Create with associations
    user := User{
        Name:  "Alice",
        Email: "alice@example.com",
    }
    
    post := Post{
        Title:   "First Post",
        Content: "Content here",
        User:    user,
        Tags: []Tag{
            {Name: "Go"},
            {Name: "Tutorial"},
        },
    }
    db.Create(&post)
    
    // 预加载 / Preload
    var posts []Post
    
    // 预加载单个关联 / Preload single association
    db.Preload("User").Find(&posts)
    
    // 预加载多个关联 / Preload multiple associations
    db.Preload("User").Preload("Tags").Find(&posts)
    
    // 嵌套预加载 / Nested preload
    db.Preload("Comments.User").Find(&posts)
    
    // 条件预加载 / Conditional preload
    db.Preload("Comments", "id > ?", 10).Find(&posts)
    
    // 自定义预加载 / Custom preload
    db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
        return db.Order("created_at DESC").Limit(5)
    }).Find(&posts)
    
    // 使用 Joins 预加载 / Preload with Joins
    db.Joins("User").Find(&posts)
    
    // 关联模式 / Association mode
    var user2 User
    db.First(&user2, 1)
    
    // 查找关联 / Find associations
    var userPosts []Post
    db.Model(&user2).Association("Posts").Find(&userPosts)
    
    // 添加关联 / Append associations
    db.Model(&user2).Association("Posts").Append(&Post{Title: "New Post"})
    
    // 替换关联 / Replace associations
    db.Model(&user2).Association("Tags").Replace(&Tag{Name: "NewTag"})
    
    // 删除关联 / Delete associations
    db.Model(&user2).Association("Posts").Delete(&post)
    
    // 清空关联 / Clear associations
    db.Model(&user2).Association("Posts").Clear()
    
    // 计数 / Count
    count := db.Model(&user2).Association("Posts").Count()
    _ = count
}
```

## 5. 事务 / Transactions

```go
package main

import (
    "errors"
    
    "gorm.io/gorm"
)

func transactionExamples(db *gorm.DB) {
    // 手动事务 / Manual transaction
    tx := db.Begin()
    
    if err := tx.Create(&User{Name: "Alice"}).Error; err != nil {
        tx.Rollback()
        return
    }
    
    if err := tx.Create(&User{Name: "Bob"}).Error; err != nil {
        tx.Rollback()
        return
    }
    
    tx.Commit()
    
    // 自动事务 / Auto transaction
    err := db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&User{Name: "Charlie"}).Error; err != nil {
            return err
        }
        
        if err := tx.Create(&User{Name: "Diana"}).Error; err != nil {
            return err
        }
        
        // 返回 nil 提交事务 / Return nil to commit
        return nil
    })
    
    if err != nil {
        // 处理错误 / Handle error
    }
    
    // 嵌套事务 (使用 SavePoint) / Nested transaction
    db.Transaction(func(tx *gorm.DB) error {
        tx.Create(&User{Name: "User1"})
        
        tx.Transaction(func(tx2 *gorm.DB) error {
            tx2.Create(&User{Name: "User2"})
            return errors.New("rollback inner")  // 回滚内部事务
        })
        
        tx.Create(&User{Name: "User3"})
        return nil  // User1 和 User3 会被保存
    })
}
```

## 6. Hooks 钩子 / Hooks

```go
package main

import (
    "fmt"
    
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type Account struct {
    ID       uint
    Username string
    Password string
}

// BeforeCreate 创建前 / Before create
func (a *Account) BeforeCreate(tx *gorm.DB) error {
    // 密码加密 / Encrypt password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    a.Password = string(hashedPassword)
    return nil
}

// AfterCreate 创建后 / After create
func (a *Account) AfterCreate(tx *gorm.DB) error {
    fmt.Printf("Account created: %d\n", a.ID)
    return nil
}

// BeforeUpdate 更新前 / Before update
func (a *Account) BeforeUpdate(tx *gorm.DB) error {
    // 如果密码被修改，重新加密
    if tx.Statement.Changed("Password") {
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
        tx.Statement.SetColumn("Password", string(hashedPassword))
    }
    return nil
}

// AfterUpdate 更新后 / After update
func (a *Account) AfterUpdate(tx *gorm.DB) error {
    fmt.Printf("Account updated: %d\n", a.ID)
    return nil
}

// BeforeDelete 删除前 / Before delete
func (a *Account) BeforeDelete(tx *gorm.DB) error {
    fmt.Printf("About to delete account: %d\n", a.ID)
    return nil
}

// AfterDelete 删除后 / After delete
func (a *Account) AfterDelete(tx *gorm.DB) error {
    fmt.Printf("Account deleted: %d\n", a.ID)
    return nil
}

// AfterFind 查询后 / After find
func (a *Account) AfterFind(tx *gorm.DB) error {
    fmt.Printf("Account found: %d\n", a.ID)
    return nil
}
```

## 7. 高级查询 / Advanced Queries

```go
package main

import (
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

func advancedQueryExamples(db *gorm.DB) {
    var users []User
    var user User
    
    // 原生 SQL / Raw SQL
    db.Raw("SELECT * FROM users WHERE age > ?", 20).Scan(&users)
    db.Exec("UPDATE users SET age = ? WHERE name = ?", 30, "Alice")
    
    // 子查询 / Subquery
    subQuery := db.Select("AVG(age)").Table("users")
    db.Where("age > (?)", subQuery).Find(&users)
    
    // 命名参数 / Named arguments
    db.Where("name = @name OR age = @age", map[string]interface{}{
        "name": "Alice",
        "age":  25,
    }).Find(&users)
    
    // Scopes - 可复用查询 / Reusable queries
    ageGreaterThan := func(age int) func(*gorm.DB) *gorm.DB {
        return func(db *gorm.DB) *gorm.DB {
            return db.Where("age > ?", age)
        }
    }
    
    active := func(db *gorm.DB) *gorm.DB {
        return db.Where("active = ?", true)
    }
    
    db.Scopes(ageGreaterThan(20), active).Find(&users)
    
    // Locking / 锁
    db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, 1)
    
    // Upsert / 插入或更新
    db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "email"}},
        DoUpdates: clause.AssignmentColumns([]string{"name", "age"}),
    }).Create(&user)
    
    // 批量插入优化 / Batch insert optimization
    var bulkUsers []User
    for i := 0; i < 1000; i++ {
        bulkUsers = append(bulkUsers, User{Name: "User"})
    }
    db.CreateInBatches(bulkUsers, 100)  // 每批 100 条
    
    // 使用 Map 结果 / Use Map for results
    var result []map[string]interface{}
    db.Model(&User{}).Find(&result)
    
    // 智能选择字段 / Smart select
    type APIUser struct {
        ID   uint
        Name string
    }
    var apiUsers []APIUser
    db.Model(&User{}).Find(&apiUsers)  // 只查询 ID 和 Name
    
    // FindInBatches - 分批处理 / Process in batches
    db.Where("age > ?", 20).FindInBatches(&users, 100, func(tx *gorm.DB, batch int) error {
        for _, user := range users {
            // 处理每条记录
            _ = user
        }
        return nil
    })
}
```

## 8. 性能优化 / Performance Optimization

```go
package main

import (
    "gorm.io/gorm"
)

func performanceExamples(db *gorm.DB) {
    // 1. 使用 Select 只查询需要的字段
    var users []User
    db.Select("id", "name").Find(&users)
    
    // 2. 使用 Joins 而不是 Preload (某些情况下更快)
    var posts []Post
    db.Joins("User").Find(&posts)
    
    // 3. 批量操作
    db.CreateInBatches(users, 100)
    
    // 4. 使用索引 (在模型定义中)
    // `gorm:"index"`
    // `gorm:"uniqueIndex"`
    
    // 5. 禁用默认事务 (对于不需要事务的操作)
    db2 := db.Session(&gorm.Session{SkipDefaultTransaction: true})
    _ = db2
    
    // 6. 准备语句缓存
    db3 := db.Session(&gorm.Session{PrepareStmt: true})
    _ = db3
    
    // 7. 使用 DryRun 调试
    stmt := db.Session(&gorm.Session{DryRun: true}).First(&User{}, 1).Statement
    _ = stmt.SQL.String()
    
    // 8. 使用连接池 (在连接时配置)
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
}
```
