# GORM 详解与实战

## 目录
- [简介](#简介)
- [安装与配置](#安装与配置)
- [数据库连接](#数据库连接)
- [模型定义](#模型定义)
- [CRUD 操作](#crud-操作)
- [高级查询](#高级查询)
- [关联关系](#关联关系)
- [事务处理](#事务处理)
- [钩子函数](#钩子函数)
- [性能优化](#性能优化)
- [实战案例](#实战案例)

---

## 简介

GORM 是 Go 语言最流行的 ORM（Object Relational Mapping）库，提供了完整的数据库操作功能。

### 主要特性
- 全功能 ORM
- 关联（Has One、Has Many、Belongs To、Many To Many、多态关联）
- 钩子（Before/After Create/Save/Update/Delete/Find）
- 预加载（Eager Loading）
- 事务、嵌套事务、保存点
- Context、预编译模式、DryRun 模式
- 批量插入、FindInBatches、FindOrCreate
- SQL Builder、Upsert、锁定、优化器/索引/注释提示
- 复合主键、索引、约束
- 自动迁移
- 自定义 Logger
- 灵活的插件 API
- 每个功能都有测试覆盖

---

## 安装与配置

### 安装 GORM

```bash
# 安装 GORM
go get -u gorm.io/gorm

# 安装数据库驱动
# MySQL
go get -u gorm.io/driver/mysql

# PostgreSQL
go get -u gorm.io/driver/postgres

# SQLite
go get -u gorm.io/driver/sqlite

# SQL Server
go get -u gorm.io/driver/sqlserver
```

---

## 数据库连接

### MySQL 连接

```go
package main

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

func ConnectMySQL() (*gorm.DB, error) {
    // DSN 格式: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
    dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    
    // 自定义日志配置
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold:             time.Second,   // 慢 SQL 阈值
            LogLevel:                  logger.Info,   // Log level
            IgnoreRecordNotFoundError: true,          // 忽略 ErrRecordNotFound 错误
            Colorful:                  true,          // 彩色打印
        },
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
        NowFunc: func() time.Time {
            return time.Now().Local()
        },
        DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
    })
    
    if err != nil {
        return nil, err
    }
    
    // 获取底层 sql.DB 对象进行连接池配置
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // 设置连接池参数
    sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
    sqlDB.SetMaxOpenConns(100)          // 最大连接数
    sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
    
    return db, nil
}
```

### PostgreSQL 连接

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectPostgreSQL() (*gorm.DB, error) {
    dsn := "host=localhost user=postgres password=password dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    return db, err
}
```

### SQLite 连接

```go
import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func ConnectSQLite() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    return db, err
}
```

---

## 模型定义

### 基本模型

```go
package models

import (
    "time"
    "gorm.io/gorm"
)

// 基础模型（推荐所有模型都嵌入）
type BaseModel struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 软删除
}

// 用户模型
type User struct {
    BaseModel
    Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
    Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
    Password  string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 不序列化
    Age       int       `gorm:"type:int;default:0" json:"age"`
    Birthday  time.Time `json:"birthday"`
    Active    bool      `gorm:"default:true" json:"active"`
    Profile   *Profile  `gorm:"foreignKey:UserID" json:"profile,omitempty"` // 一对一
    Posts     []Post    `gorm:"foreignKey:AuthorID" json:"posts,omitempty"` // 一对多
}

// 用户资料
type Profile struct {
    BaseModel
    UserID   uint   `gorm:"uniqueIndex;not null" json:"user_id"`
    Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
    Bio      string `gorm:"type:text" json:"bio"`
    Location string `gorm:"type:varchar(100)" json:"location"`
}

// 文章模型
type Post struct {
    BaseModel
    Title    string    `gorm:"type:varchar(200);not null;index" json:"title"`
    Content  string    `gorm:"type:text" json:"content"`
    AuthorID uint      `gorm:"not null;index" json:"author_id"`
    Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
    Tags     []Tag     `gorm:"many2many:post_tags;" json:"tags,omitempty"` // 多对多
    ViewCount int      `gorm:"default:0" json:"view_count"`
    Published bool     `gorm:"default:false" json:"published"`
}

// 标签模型
type Tag struct {
    BaseModel
    Name  string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
    Posts []Post `gorm:"many2many:post_tags;" json:"posts,omitempty"`
}

// 表名自定义
func (User) TableName() string {
    return "users"
}

func (Profile) TableName() string {
    return "user_profiles"
}

func (Post) TableName() string {
    return "posts"
}

func (Tag) TableName() string {
    return "tags"
}
```

### 字段标签详解

```go
type Example struct {
    // 主键
    ID uint `gorm:"primaryKey"`
    
    // 列名
    Name string `gorm:"column:user_name"`
    
    // 数据类型
    Age int `gorm:"type:int"`
    
    // 大小
    Code string `gorm:"size:100"`
    
    // 精度
    Price float64 `gorm:"precision:10;scale:2"`
    
    // 非空
    Email string `gorm:"not null"`
    
    // 唯一
    Phone string `gorm:"unique"`
    
    // 默认值
    Status string `gorm:"default:active"`
    
    // 索引
    Title string `gorm:"index"`
    
    // 唯一索引
    Username string `gorm:"uniqueIndex"`
    
    // 复合索引
    GroupID uint   `gorm:"index:idx_group_user"`
    UserID  uint   `gorm:"index:idx_group_user"`
    
    // 忽略字段
    TempField string `gorm:"-"`
    
    // 只读字段
    ReadOnly string `gorm:"<-:false"`
    
    // 只写字段
    WriteOnly string `gorm:"->:false"`
    
    // 只在创建时写入
    CreateOnly string `gorm:"<-:create"`
    
    // 自动时间戳
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    
    // 软删除
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 自动迁移

```go
func AutoMigrate(db *gorm.DB) error {
    // 自动迁移会创建表、缺失的外键、约束、列和索引
    // 但不会改变现有列的类型或删除未使用的列
    return db.AutoMigrate(
        &User{},
        &Profile{},
        &Post{},
        &Tag{},
    )
}
```

---

## CRUD 操作

### 创建（Create）

```go
// 1. 创建单条记录
func CreateUser(db *gorm.DB) {
    user := User{
        Username: "zhangsan",
        Email:    "zhangsan@example.com",
        Password: "hashed_password",
        Age:      25,
        Birthday: time.Now().AddDate(-25, 0, 0),
    }
    
    result := db.Create(&user)
    if result.Error != nil {
        log.Printf("创建失败: %v", result.Error)
        return
    }
    
    log.Printf("创建成功，ID: %d, 影响行数: %d", user.ID, result.RowsAffected)
}

// 2. 批量创建
func BatchCreateUsers(db *gorm.DB) {
    users := []User{
        {Username: "user1", Email: "user1@example.com", Password: "pass1"},
        {Username: "user2", Email: "user2@example.com", Password: "pass2"},
        {Username: "user3", Email: "user3@example.com", Password: "pass3"},
    }
    
    // 批量插入，单次插入 100 条
    result := db.CreateInBatches(users, 100)
    log.Printf("批量创建 %d 条记录", result.RowsAffected)
}

// 3. 选定字段创建
func CreateWithSelectedFields(db *gorm.DB) {
    user := User{
        Username: "lisi",
        Email:    "lisi@example.com",
        Password: "password",
        Age:      30,
    }
    
    // 只插入指定字段
    db.Select("Username", "Email", "Password").Create(&user)
    
    // 或者排除某些字段
    db.Omit("Age").Create(&user)
}

// 4. 使用 Map 创建
func CreateFromMap(db *gorm.DB) {
    db.Model(&User{}).Create(map[string]interface{}{
        "Username": "wangwu",
        "Email":    "wangwu@example.com",
        "Password": "password",
    })
}

// 5. 冲突时更新（Upsert）
func UpsertUser(db *gorm.DB) {
    user := User{
        Username: "zhangsan",
        Email:    "zhangsan@example.com",
        Age:      26,
    }
    
    // 当 email 冲突时更新 age 和 updated_at
    db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "email"}},
        DoUpdates: clause.AssignmentColumns([]string{"age", "updated_at"}),
    }).Create(&user)
}
```

### 查询（Read）

```go
// 1. 查询单条记录
func GetUser(db *gorm.DB, id uint) (*User, error) {
    var user User
    
    // First 方法查询第一条记录
    result := db.First(&user, id)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("用户不存在")
        }
        return nil, result.Error
    }
    
    return &user, nil
}

// 2. 条件查询
func FindUsers(db *gorm.DB) {
    var users []User
    
    // Where 条件查询
    db.Where("age > ?", 18).Find(&users)
    
    // 多个条件
    db.Where("age > ? AND active = ?", 18, true).Find(&users)
    
    // 或条件
    db.Where("username = ? OR email = ?", "zhangsan", "test@example.com").Find(&users)
    
    // IN 查询
    db.Where("username IN ?", []string{"zhangsan", "lisi", "wangwu"}).Find(&users)
    
    // LIKE 查询
    db.Where("username LIKE ?", "%zhang%").Find(&users)
    
    // BETWEEN 查询
    db.Where("age BETWEEN ? AND ?", 18, 30).Find(&users)
    
    // Struct 条件（只会查询非零值字段）
    db.Where(&User{Username: "zhangsan", Active: true}).Find(&users)
    
    // Map 条件
    db.Where(map[string]interface{}{"username": "zhangsan", "age": 25}).Find(&users)
}

// 3. 排序、限制、偏移
func QueryWithOptions(db *gorm.DB) {
    var users []User
    
    // 排序
    db.Order("age desc, created_at asc").Find(&users)
    
    // 限制数量
    db.Limit(10).Find(&users)
    
    // 偏移
    db.Offset(5).Limit(10).Find(&users)
    
    // 组合使用
    db.Where("age > ?", 18).
       Order("created_at desc").
       Limit(20).
       Offset(0).
       Find(&users)
}

// 4. 选择特定字段
func SelectFields(db *gorm.DB) {
    var users []User
    
    // 选择特定字段
    db.Select("id", "username", "email").Find(&users)
    
    // 排除字段
    db.Omit("password").Find(&users)
    
    // 查询到 map
    var result []map[string]interface{}
    db.Model(&User{}).Select("id", "username").Find(&result)
}

// 5. 分组和聚合
func GroupAndAggregate(db *gorm.DB) {
    // 统计数量
    var count int64
    db.Model(&User{}).Where("age > ?", 18).Count(&count)
    
    // 分组统计
    type Result struct {
        Age   int
        Count int
    }
    var results []Result
    db.Model(&User{}).
       Select("age, count(*) as count").
       Group("age").
       Having("count > ?", 1).
       Find(&results)
    
    // 聚合函数
    var avgAge float64
    db.Model(&User{}).Select("AVG(age)").Row().Scan(&avgAge)
}

// 6. 子查询
func SubQuery(db *gorm.DB) {
    var users []User
    
    // 子查询
    subQuery := db.Model(&Post{}).Select("author_id").Where("published = ?", true)
    db.Where("id IN (?)", subQuery).Find(&users)
}

// 7. 原生 SQL
func RawSQL(db *gorm.DB) {
    var users []User
    
    // 原生 SQL 查询
    db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
    
    // 执行原生 SQL
    db.Exec("UPDATE users SET active = ? WHERE age < ?", false, 18)
}
```

### 更新（Update）

```go
// 1. 更新单个字段
func UpdateSingleField(db *gorm.DB, userID uint) {
    // 更新单个字段
    db.Model(&User{}).Where("id = ?", userID).Update("age", 26)
    
    // 使用 struct 更新（只更新非零值字段）
    db.Model(&User{}).Where("id = ?", userID).Updates(&User{Age: 27, Active: true})
}

// 2. 更新多个字段
func UpdateMultipleFields(db *gorm.DB, userID uint) {
    // 使用 map 更新（可以更新零值）
    db.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
        "age":    0,
        "active": false,
    })
    
    // 使用 Select 选择要更新的字段
    user := User{Age: 28, Active: false}
    db.Model(&User{}).Where("id = ?", userID).Select("Age", "Active").Updates(&user)
    
    // 使用 Omit 排除某些字段
    db.Model(&User{}).Where("id = ?", userID).Omit("CreatedAt").Updates(&user)
}

// 3. 批量更新
func BatchUpdate(db *gorm.DB) {
    // 批量更新
    db.Model(&User{}).Where("age < ?", 18).Updates(map[string]interface{}{
        "active": false,
    })
}

// 4. 更新表达式
func UpdateWithExpression(db *gorm.DB) {
    // 使用表达式更新
    db.Model(&Post{}).Where("id = ?", 1).Update("view_count", gorm.Expr("view_count + ?", 1))
    
    // 使用子查询更新
    db.Model(&User{}).Where("id = ?", 1).Update("age", gorm.Expr(
        "(SELECT AVG(age) FROM users WHERE active = ?)", true,
    ))
}

// 5. 返回更新后的数据
func UpdateAndReturn(db *gorm.DB, userID uint) (*User, error) {
    var user User
    result := db.Model(&User{}).
        Clauses(clause.Returning{}).
        Where("id = ?", userID).
        Updates(map[string]interface{}{"age": 30})
    
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &user, nil
}
```

### 删除（Delete）

```go
// 1. 软删除（推荐）
func SoftDelete(db *gorm.DB, userID uint) {
    // 软删除（设置 deleted_at 字段）
    result := db.Delete(&User{}, userID)
    log.Printf("软删除 %d 条记录", result.RowsAffected)
    
    // 条件软删除
    db.Where("age < ?", 18).Delete(&User{})
}

// 2. 永久删除
func HardDelete(db *gorm.DB, userID uint) {
    // 永久删除（真正从数据库删除）
    db.Unscoped().Delete(&User{}, userID)
}

// 3. 查询包含软删除的记录
func QueryWithDeleted(db *gorm.DB) {
    var users []User
    
    // 包含软删除的记录
    db.Unscoped().Find(&users)
    
    // 只查询软删除的记录
    db.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
}

// 4. 批量删除
func BatchDelete(db *gorm.DB) {
    // 批量软删除
    db.Where("age < ?", 18).Delete(&User{})
    
    // 批量硬删除
    db.Unscoped().Where("age < ?", 18).Delete(&User{})
}
```

---

## 高级查询

### 预加载（Eager Loading）

```go
// 1. 预加载关联数据
func PreloadAssociations(db *gorm.DB) {
    var users []User
    
    // 预加载 Profile
    db.Preload("Profile").Find(&users)
    
    // 预加载多个关联
    db.Preload("Profile").Preload("Posts").Find(&users)
    
    // 预加载所有关联
    db.Preload(clause.Associations).Find(&users)
    
    // 嵌套预加载
    db.Preload("Posts.Tags").Find(&users)
    
    // 条件预加载
    db.Preload("Posts", "published = ?", true).Find(&users)
    
    // 自定义预加载 SQL
    db.Preload("Posts", func(db *gorm.DB) *gorm.DB {
        return db.Order("posts.created_at DESC").Limit(5)
    }).Find(&users)
}

// 2. Joins 预加载
func JoinsPreload(db *gorm.DB) {
    var users []User
    
    // Joins 预加载（使用 LEFT JOIN）
    db.Joins("Profile").Find(&users)
    
    // 带条件的 Joins
    db.Joins("Profile").Where("Profile.location = ?", "Beijing").Find(&users)
    
    // 多个 Joins
    db.Joins("Profile").Joins("LEFT JOIN posts ON posts.author_id = users.id").Find(&users)
}
```

### 作用域（Scopes）

```go
// 定义作用域
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func AgeGreaterThan(age int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}

func OrderByCreatedAt(db *gorm.DB) *gorm.DB {
    return db.Order("created_at desc")
}

// 使用作用域
func UseScopes(db *gorm.DB) {
    var users []User
    
    // 使用单个作用域
    db.Scopes(ActiveUsers).Find(&users)
    
    // 使用多个作用域
    db.Scopes(ActiveUsers, AgeGreaterThan(18), OrderByCreatedAt).Find(&users)
}

// 分页作用域
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        offset := (page - 1) * pageSize
        return db.Offset(offset).Limit(pageSize)
    }
}

func PaginateUsers(db *gorm.DB, page, pageSize int) {
    var users []User
    db.Scopes(Paginate(page, pageSize)).Find(&users)
}
```

### 查询链

```go
func ChainQuery(db *gorm.DB) {
    // 创建查询链
    query := db.Model(&User{})
    
    // 动态添加条件
    if username != "" {
        query = query.Where("username LIKE ?", "%"+username+"%")
    }
    
    if minAge > 0 {
        query = query.Where("age >= ?", minAge)
    }
    
    if maxAge > 0 {
        query = query.Where("age <= ?", maxAge)
    }
    
    // 执行查询
    var users []User
    query.Find(&users)
}
```

---

## 关联关系

### 一对一（Has One）

```go
// 创建带关联的记录
func CreateUserWithProfile(db *gorm.DB) {
    user := User{
        Username: "zhangsan",
        Email:    "zhangsan@example.com",
        Password: "password",
        Profile: &Profile{
            Avatar:   "avatar.jpg",
            Bio:      "This is my bio",
            Location: "Beijing",
        },
    }
    
    // 会自动创建 user 和 profile
    db.Create(&user)
}

// 关联查询
func QueryUserWithProfile(db *gorm.DB, userID uint) {
    var user User
    db.Preload("Profile").First(&user, userID)
}

// 更新关联
func UpdateProfile(db *gorm.DB, userID uint) {
    var user User
    db.First(&user, userID)
    
    // 更新关联
    db.Model(&user).Association("Profile").Replace(&Profile{
        Avatar:   "new_avatar.jpg",
        Bio:      "Updated bio",
        Location: "Shanghai",
    })
}
```

### 一对多（Has Many）

```go
// 创建带文章的用户
func CreateUserWithPosts(db *gorm.DB) {
    user := User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password",
        Posts: []Post{
            {Title: "First Post", Content: "Content 1", Published: true},
            {Title: "Second Post", Content: "Content 2", Published: false},
        },
    }
    
    db.Create(&user)
}

// 查询用户及其文章
func QueryUserWithPosts(db *gorm.DB, userID uint) {
    var user User
    db.Preload("Posts").First(&user, userID)
    
    // 带条件的预加载
    db.Preload("Posts", "published = ?", true).First(&user, userID)
}

// 添加关联
func AddPostToUser(db *gorm.DB, userID uint) {
    var user User
    db.First(&user, userID)
    
    post := Post{
        Title:   "New Post",
        Content: "New Content",
    }
    
    // 添加关联
    db.Model(&user).Association("Posts").Append(&post)
}

// 删除关联
func RemovePostFromUser(db *gorm.DB, userID, postID uint) {
    var user User
    db.First(&user, userID)
    
    var post Post
    db.First(&post, postID)
    
    // 删除关联（不删除 post 记录）
    db.Model(&user).Association("Posts").Delete(&post)
}

// 清空所有关联
func ClearUserPosts(db *gorm.DB, userID uint) {
    var user User
    db.First(&user, userID)
    
    db.Model(&user).Association("Posts").Clear()
}

// 统计关联数量
func CountUserPosts(db *gorm.DB, userID uint) int64 {
    var user User
    db.First(&user, userID)
    
    return db.Model(&user).Association("Posts").Count()
}
```

### 多对多（Many To Many）

```go
// 创建带标签的文章
func CreatePostWithTags(db *gorm.DB) {
    post := Post{
        Title:   "Article Title",
        Content: "Article Content",
        Tags: []Tag{
            {Name: "Go"},
            {Name: "Programming"},
            {Name: "Backend"},
        },
    }
    
    db.Create(&post)
}

// 查询文章及其标签
func QueryPostWithTags(db *gorm.DB, postID uint) {
    var post Post
    db.Preload("Tags").First(&post, postID)
}

// 添加标签到文章
func AddTagsToPost(db *gorm.DB, postID uint) {
    var post Post
    db.First(&post, postID)
    
    tags := []Tag{
        {Name: "Web"},
        {Name: "API"},
    }
    
    // 添加多对多关联
    db.Model(&post).Association("Tags").Append(&tags)
}

// 替换所有标签
func ReplacePostTags(db *gorm.DB, postID uint) {
    var post Post
    db.First(&post, postID)
    
    newTags := []Tag{
        {Name: "Updated"},
        {Name: "Tags"},
    }
    
    db.Model(&post).Association("Tags").Replace(&newTags)
}

// 查询包含特定标签的文章
func FindPostsByTag(db *gorm.DB, tagName string) {
    var posts []Post
    db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
       Joins("JOIN tags ON tags.id = post_tags.tag_id").
       Where("tags.name = ?", tagName).
       Find(&posts)
}
```

---

## 事务处理

### 手动事务

```go
// 1. 基本事务
func TransferMoney(db *gorm.DB, fromUserID, toUserID uint, amount float64) error {
    // 开始事务
    tx := db.Begin()
    
    // 发生错误时回滚
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // 检查事务是否开启成功
    if tx.Error != nil {
        return tx.Error
    }
    
    // 扣除转出方余额
    if err := tx.Model(&Account{}).Where("user_id = ?", fromUserID).
        Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 增加转入方余额
    if err := tx.Model(&Account{}).Where("user_id = ?", toUserID).
        Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
        tx.Rollback()
        return err
    }
    
    // 提交事务
    return tx.Commit().Error
}
```

### 事务回调

```go
// 2. 使用 Transaction 方法（推荐）
func CreateUserWithProfileTx(db *gorm.DB, user *User) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        if err := tx.Create(user).Error; err != nil {
            return err // 返回错误会自动回滚
        }
        
        // 创建用户资料
        profile := Profile{
            UserID:   user.ID,
            Avatar:   "default.jpg",
            Location: "Unknown",
        }
        if err := tx.Create(&profile).Error; err != nil {
            return err
        }
        
        // 返回 nil 提交事务
        return nil
    })
}
```

### 嵌套事务

```go
// 3. 嵌套事务
func NestedTransaction(db *gorm.DB) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 外层事务操作
        if err := tx.Create(&User{Username: "user1"}).Error; err != nil {
            return err
        }
        
        // 嵌套事务
        return tx.Transaction(func(tx2 *gorm.DB) error {
            // 内层事务操作
            if err := tx2.Create(&User{Username: "user2"}).Error; err != nil {
                return err
            }
            return nil
        })
    })
}
```

### 保存点

```go
// 4. 保存点
func SavepointExample(db *gorm.DB) error {
    tx := db.Begin()
    
    tx.Create(&User{Username: "user1"})
    
    // 创建保存点
    tx.SavePoint("sp1")
    
    tx.Create(&User{Username: "user2"})
    
    // 回滚到保存点
    tx.RollbackTo("sp1")
    
    // 提交事务（只有 user1 被创建）
    return tx.Commit().Error
}
```

---

## 钩子函数

### Before/After 钩子

```go
type User struct {
    BaseModel
    Username string
    Password string
    Email    string
}

// BeforeSave - 保存之前
func (u *User) BeforeSave(tx *gorm.DB) error {
    // 可以在这里做数据验证
    if u.Username == "" {
        return errors.New("username cannot be empty")
    }
    return nil
}

// BeforeCreate - 创建之前
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 例如：加密密码
    if u.Password != "" {
        hashedPassword, err := hashPassword(u.Password)
        if err != nil {
            return err
        }
        u.Password = hashedPassword
    }
    return nil
}

// AfterCreate - 创建之后
func (u *User) AfterCreate(tx *gorm.DB) error {
    // 例如：发送欢迎邮件
    log.Printf("New user created: %s", u.Username)
    return nil
}

// BeforeUpdate - 更新之前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // 检查是否有权限更新
    return nil
}

// AfterUpdate - 更新之后
func (u *User) AfterUpdate(tx *gorm.DB) error {
    log.Printf("User updated: %s", u.Username)
    return nil
}

// BeforeDelete - 删除之前
func (u *User) BeforeDelete(tx *gorm.DB) error {
    // 检查是否可以删除
    return nil
}

// AfterDelete - 删除之后
func (u *User) AfterDelete(tx *gorm.DB) error {
    log.Printf("User deleted: %s", u.Username)
    return nil
}

// AfterFind - 查询之后
func (u *User) AfterFind(tx *gorm.DB) error {
    // 可以在这里做数据处理
    return nil
}

func hashPassword(password string) (string, error) {
    // 实现密码加密逻辑
    return "hashed_" + password, nil
}
```

---

## 性能优化

### 索引优化

```go
type User struct {
    BaseModel
    Username string `gorm:"index:idx_username"`           // 单字段索引
    Email    string `gorm:"uniqueIndex"`                  // 唯一索引
    Age      int    `gorm:"index:idx_age_active"`         // 复合索引
    Active   bool   `gorm:"index:idx_age_active"`         // 复合索引
    Phone    string `gorm:"index:,type:btree,length:10"`  // 指定索引类型和长度
}
```

### 批量操作

```go
// 1. 批量插入
func BatchInsert(db *gorm.DB) {
    users := make([]User, 1000)
    for i := 0; i < 1000; i++ {
        users[i] = User{
            Username: fmt.Sprintf("user%d", i),
            Email:    fmt.Sprintf("user%d@example.com", i),
        }
    }
    
    // 每次插入 100 条
    db.CreateInBatches(users, 100)
}

// 2. 批量查询
func FindInBatches(db *gorm.DB) {
    result := db.Where("active = ?", true).FindInBatches(&[]User{}, 100, func(tx *gorm.DB, batch int) error {
        // 处理每批数据
        for _, user := range *tx.Statement.Dest.(*[]User) {
            // 处理用户
            processUser(user)
        }
        return nil
    })
    
    if result.Error != nil {
        log.Printf("Error: %v", result.Error)
    }
}

func processUser(user User) {
    // 处理用户逻辑
}
```

### 预编译

```go
// 预编译 SQL
func PreparedStatement(db *gorm.DB) {
    // 开启预编译模式
    stmt := db.Session(&gorm.Session{PrepareStmt: true})
    
    // 多次使用预编译语句
    for i := 0; i < 100; i++ {
        var user User
        stmt.First(&user, i)
    }
}
```

### 选择性查询

```go
// 只查询需要的字段
func SelectiveQuery(db *gorm.DB) {
    var users []User
    
    // 只查询 ID、Username 和 Email
    db.Select("id", "username", "email").Find(&users)
    
    // 使用 struct 只查询某些字段
    type UserBasic struct {
        ID       uint
        Username string
        Email    string
    }
    var userBasics []UserBasic
    db.Model(&User{}).Find(&userBasics)
}
```

### 连接池配置

```go
func ConfigureConnectionPool(db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        panic(err)
    }
    
    // 设置连接池参数
    sqlDB.SetMaxIdleConns(10)                // 最大空闲连接数
    sqlDB.SetMaxOpenConns(100)               // 最大打开连接数
    sqlDB.SetConnMaxLifetime(time.Hour)      // 连接最大生存时间
    sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 连接最大空闲时间
}
```

---

## 实战案例

### 完整的用户管理系统

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// 数据库管理器
type DBManager struct {
    db *gorm.DB
}

// 初始化数据库
func NewDBManager(dsn string) (*DBManager, error) {
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().Local()
        },
    })
    if err != nil {
        return nil, err
    }
    
    // 自动迁移
    if err := db.AutoMigrate(&User{}, &Profile{}, &Post{}, &Tag{}); err != nil {
        return nil, err
    }
    
    return &DBManager{db: db}, nil
}

// 用户服务
type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// 创建用户
func (s *UserService) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
    user := &User{
        Username: username,
        Email:    email,
        Password: password,
    }
    
    if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    return user, nil
}

// 根据 ID 获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*User, error) {
    var user User
    err := s.db.WithContext(ctx).
        Preload("Profile").
        Preload("Posts").
        First(&user, id).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }
    
    return &user, nil
}

// 更新用户
func (s *UserService) UpdateUser(ctx context.Context, id uint, updates map[string]interface{}) error {
    result := s.db.WithContext(ctx).
        Model(&User{}).
        Where("id = ?", id).
        Updates(updates)
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

// 删除用户（软删除）
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
    result := s.db.WithContext(ctx).Delete(&User{}, id)
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("user not found")
    }
    
    return nil
}

// 分页查询用户
type PaginationParams struct {
    Page     int
    PageSize int
}

type PaginationResult struct {
    Total    int64
    Page     int
    PageSize int
    Data     interface{}
}

func (s *UserService) ListUsers(ctx context.Context, params PaginationParams) (*PaginationResult, error) {
    var users []User
    var total int64
    
    // 计算总数
    if err := s.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
        return nil, err
    }
    
    // 分页查询
    offset := (params.Page - 1) * params.PageSize
    if err := s.db.WithContext(ctx).
        Offset(offset).
        Limit(params.PageSize).
        Order("created_at DESC").
        Find(&users).Error; err != nil {
        return nil, err
    }
    
    return &PaginationResult{
        Total:    total,
        Page:     params.Page,
        PageSize: params.PageSize,
        Data:     users,
    }, nil
}

// 搜索用户
func (s *UserService) SearchUsers(ctx context.Context, keyword string) ([]User, error) {
    var users []User
    err := s.db.WithContext(ctx).
        Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
        Find(&users).Error
    
    return users, err
}

// 文章服务
type PostService struct {
    db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
    return &PostService{db: db}
}

// 创建文章
func (s *PostService) CreatePost(ctx context.Context, authorID uint, title, content string, tagNames []string) (*Post, error) {
    post := &Post{
        Title:    title,
        Content:  content,
        AuthorID: authorID,
    }
    
    // 使用事务
    err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 创建文章
        if err := tx.Create(post).Error; err != nil {
            return err
        }
        
        // 处理标签
        if len(tagNames) > 0 {
            var tags []Tag
            for _, name := range tagNames {
                var tag Tag
                // 查找或创建标签
                if err := tx.Where("name = ?", name).FirstOrCreate(&tag, Tag{Name: name}).Error; err != nil {
                    return err
                }
                tags = append(tags, tag)
            }
            
            // 关联标签
            if err := tx.Model(post).Association("Tags").Append(&tags); err != nil {
                return err
            }
        }
        
        return nil
    })
    
    if err != nil {
        return nil, err
    }
    
    return post, nil
}

// 获取文章详情
func (s *PostService) GetPostByID(ctx context.Context, id uint) (*Post, error) {
    var post Post
    err := s.db.WithContext(ctx).
        Preload("Author").
        Preload("Tags").
        First(&post, id).Error
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("post not found")
        }
        return nil, err
    }
    
    return &post, nil
}

// 增加浏览量
func (s *PostService) IncrementViewCount(ctx context.Context, id uint) error {
    return s.db.WithContext(ctx).
        Model(&Post{}).
        Where("id = ?", id).
        Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// 发布文章
func (s *PostService) PublishPost(ctx context.Context, id uint) error {
    return s.db.WithContext(ctx).
        Model(&Post{}).
        Where("id = ?", id).
        Update("published", true).Error
}

// 获取用户的所有文章
func (s *PostService) GetUserPosts(ctx context.Context, userID uint) ([]Post, error) {
    var posts []Post
    err := s.db.WithContext(ctx).
        Where("author_id = ?", userID).
        Preload("Tags").
        Order("created_at DESC").
        Find(&posts).Error
    
    return posts, err
}

// 主函数示例
func main() {
    // 初始化数据库
    dsn := "root:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
    dbManager, err := NewDBManager(dsn)
    if err != nil {
        panic(err)
    }
    
    // 创建服务
    userService := NewUserService(dbManager.db)
    postService := NewPostService(dbManager.db)
    
    ctx := context.Background()
    
    // 创建用户
    user, err := userService.CreateUser(ctx, "zhangsan", "zhangsan@example.com", "password123")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created user: %+v\n", user)
    
    // 创建文章
    post, err := postService.CreatePost(ctx, user.ID, "My First Post", "This is the content", []string{"Go", "GORM"})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created post: %+v\n", post)
    
    // 查询用户
    foundUser, err := userService.GetUserByID(ctx, user.ID)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found user: %+v\n", foundUser)
    
    // 分页查询
    result, err := userService.ListUsers(ctx, PaginationParams{Page: 1, PageSize: 10})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Pagination result: %+v\n", result)
}
```

---

## 总结

本文详细介绍了 GORM 的各个方面：

1. **基础功能**：数据库连接、模型定义、CRUD 操作
2. **高级功能**：关联关系、预加载、作用域、事务
3. **性能优化**：索引、批量操作、连接池配置
4. **实战应用**：完整的用户管理系统示例

GORM 是一个功能强大、易于使用的 ORM 库，掌握这些知识可以帮助你高效地开发 Go 应用程序。

### 最佳实践

1. 使用 `BaseModel` 统一管理公共字段
2. 合理使用软删除
3. 利用预加载避免 N+1 查询问题
4. 使用事务保证数据一致性
5. 合理配置索引提升查询性能
6. 使用连接池优化数据库连接
7. 使用 Context 支持超时和取消操作
8. 使用钩子函数处理业务逻辑
9. 使用作用域封装常用查询条件
10. 定期监控慢 SQL 并优化

更多详细信息请参考 [GORM 官方文档](https://gorm.io/docs/)。
