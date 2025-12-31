# MySQL 原生驱动 / MySQL Native Driver

## 1. 连接与配置 / Connection & Configuration

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    _ "github.com/go-sql-driver/mysql"
)

// 安装 / Installation:
// go get -u github.com/go-sql-driver/mysql

func main() {
    // DSN 格式 / DSN format:
    // [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...]
    
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    
    // 打开连接 / Open connection
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
    defer db.Close()
    
    // 验证连接 / Verify connection
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping database:", err)
    }
    
    // 配置连接池 / Configure connection pool
    db.SetMaxOpenConns(25)                  // 最大打开连接数
    db.SetMaxIdleConns(5)                   // 最大空闲连接数
    db.SetConnMaxLifetime(5 * time.Minute)  // 连接最大生存时间
    db.SetConnMaxIdleTime(5 * time.Minute)  // 空闲连接最大时间
    
    // 查看连接池状态 / Check pool stats
    stats := db.Stats()
    fmt.Printf("Open connections: %d\n", stats.OpenConnections)
    fmt.Printf("In use: %d\n", stats.InUse)
    fmt.Printf("Idle: %d\n", stats.Idle)
    
    fmt.Println("Database connected!")
}

// 创建表 / Create table
const createTableSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE,
    age INT DEFAULT 0,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

func createTable(db *sql.DB) error {
    _, err := db.Exec(createTableSQL)
    return err
}
```

## 2. 基本 CRUD 操作 / Basic CRUD Operations

```go
package main

import (
    "database/sql"
    "fmt"
    "time"
)

type User struct {
    ID        int
    Name      string
    Email     sql.NullString  // 可空字段
    Age       int
    Active    bool
    CreatedAt time.Time
}

// Create 创建 / Create
func createUser(db *sql.DB, name, email string, age int) (int64, error) {
    result, err := db.Exec(
        "INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
        name, email, age,
    )
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

// Read 查询单条 / Query single row
func getUserByID(db *sql.DB, id int) (*User, error) {
    user := &User{}
    err := db.QueryRow(
        "SELECT id, name, email, age, active, created_at FROM users WHERE id = ?",
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.Active, &user.CreatedAt)
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user not found")
    }
    if err != nil {
        return nil, err
    }
    return user, nil
}

// Read 查询多条 / Query multiple rows
func getUsers(db *sql.DB, minAge int) ([]User, error) {
    rows, err := db.Query(
        "SELECT id, name, email, age, active, created_at FROM users WHERE age >= ?",
        minAge,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Active, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    
    // 检查遍历过程中的错误
    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return users, nil
}

// Update 更新 / Update
func updateUser(db *sql.DB, id int, name string, age int) (int64, error) {
    result, err := db.Exec(
        "UPDATE users SET name = ?, age = ? WHERE id = ?",
        name, age, id,
    )
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}

// Delete 删除 / Delete
func deleteUser(db *sql.DB, id int) (int64, error) {
    result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
    if err != nil {
        return 0, err
    }
    return result.RowsAffected()
}

// 批量插入 / Batch insert
func batchInsertUsers(db *sql.DB, users []User) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    stmt, err := tx.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()
    
    for _, u := range users {
        _, err := stmt.Exec(u.Name, u.Email.String, u.Age)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    
    return tx.Commit()
}
```

## 3. 事务处理 / Transaction Handling

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
)

// 基本事务 / Basic transaction
func transferMoney(db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    // 使用 defer 确保事务被处理
    defer func() {
        if err != nil {
            tx.Rollback()
            return
        }
    }()
    
    // 扣款
    _, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        return err
    }
    
    // 入账
    _, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

// 带 Context 的事务 / Transaction with context
func transferMoneyWithContext(ctx context.Context, db *sql.DB, fromID, toID int, amount float64) error {
    tx, err := db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,  // 设置隔离级别
        ReadOnly:  false,
    })
    if err != nil {
        return err
    }
    
    defer tx.Rollback()  // 如果已提交，Rollback 会被忽略
    
    // 执行操作...
    _, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromID)
    if err != nil {
        return err
    }
    
    _, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toID)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

// 事务辅助函数 / Transaction helper function
func withTransaction(db *sql.DB, fn func(*sql.Tx) error) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    if err := fn(tx); err != nil {
        tx.Rollback()
        return err
    }
    
    return tx.Commit()
}

// 使用辅助函数 / Using helper function
func example(db *sql.DB) error {
    return withTransaction(db, func(tx *sql.Tx) error {
        _, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
        if err != nil {
            return err
        }
        
        _, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "Bob")
        return err
    })
}
```

## 4. 预处理语句 / Prepared Statements

```go
package main

import (
    "database/sql"
    "fmt"
)

func preparedStatementExamples(db *sql.DB) {
    // 创建预处理语句 / Create prepared statement
    stmt, err := db.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        return
    }
    defer stmt.Close()
    
    // 多次使用同一个预处理语句 / Use statement multiple times
    users := []struct {
        Name  string
        Email string
        Age   int
    }{
        {"Alice", "alice@example.com", 25},
        {"Bob", "bob@example.com", 30},
        {"Charlie", "charlie@example.com", 35},
    }
    
    for _, u := range users {
        result, err := stmt.Exec(u.Name, u.Email, u.Age)
        if err != nil {
            continue
        }
        id, _ := result.LastInsertId()
        fmt.Printf("Inserted user with ID: %d\n", id)
    }
    
    // 预处理查询语句 / Prepared query statement
    queryStmt, err := db.Prepare("SELECT id, name, age FROM users WHERE age > ?")
    if err != nil {
        return
    }
    defer queryStmt.Close()
    
    rows, err := queryStmt.Query(20)
    if err != nil {
        return
    }
    defer rows.Close()
    
    for rows.Next() {
        var id, age int
        var name string
        rows.Scan(&id, &name, &age)
        fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
    }
}
```

## 5. NULL 值处理 / NULL Value Handling

```go
package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
)

// 使用 sql.Null* 类型 / Using sql.Null* types
type UserWithNulls struct {
    ID        int
    Name      string
    Email     sql.NullString
    Age       sql.NullInt64
    Score     sql.NullFloat64
    Active    sql.NullBool
    CreatedAt sql.NullTime
}

func handleNullValues(db *sql.DB) {
    var user UserWithNulls
    
    err := db.QueryRow(`
        SELECT id, name, email, age, score, active, created_at 
        FROM users WHERE id = ?
    `, 1).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.Age,
        &user.Score,
        &user.Active,
        &user.CreatedAt,
    )
    
    if err != nil {
        return
    }
    
    // 检查是否为 NULL / Check if NULL
    if user.Email.Valid {
        fmt.Println("Email:", user.Email.String)
    } else {
        fmt.Println("Email is NULL")
    }
    
    if user.Age.Valid {
        fmt.Println("Age:", user.Age.Int64)
    } else {
        fmt.Println("Age is NULL")
    }
    
    // 插入 NULL 值 / Insert NULL value
    db.Exec(
        "INSERT INTO users (name, email) VALUES (?, ?)",
        "TestUser",
        sql.NullString{Valid: false},  // 插入 NULL
    )
    
    // 使用指针处理 NULL / Using pointers for NULL
    type UserWithPointers struct {
        ID    int
        Name  string
        Email *string  // nil 表示 NULL
        Age   *int
    }
    
    var user2 UserWithPointers
    db.QueryRow("SELECT id, name, email, age FROM users WHERE id = ?", 1).Scan(
        &user2.ID, &user2.Name, &user2.Email, &user2.Age,
    )
    
    if user2.Email != nil {
        fmt.Println("Email:", *user2.Email)
    }
}

// 自定义 NULL 类型 (支持 JSON) / Custom NULL type (with JSON support)
type NullString struct {
    sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
    if ns.Valid {
        return json.Marshal(ns.String)
    }
    return []byte("null"), nil
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
    if string(data) == "null" {
        ns.Valid = false
        return nil
    }
    ns.Valid = true
    return json.Unmarshal(data, &ns.String)
}
```

## 6. 高级查询 / Advanced Queries

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"
)

// 动态构建查询 / Dynamic query building
func searchUsers(db *sql.DB, name string, minAge, maxAge int, active *bool) ([]User, error) {
    query := "SELECT id, name, email, age, active, created_at FROM users WHERE 1=1"
    args := []interface{}{}
    
    if name != "" {
        query += " AND name LIKE ?"
        args = append(args, "%"+name+"%")
    }
    
    if minAge > 0 {
        query += " AND age >= ?"
        args = append(args, minAge)
    }
    
    if maxAge > 0 {
        query += " AND age <= ?"
        args = append(args, maxAge)
    }
    
    if active != nil {
        query += " AND active = ?"
        args = append(args, *active)
    }
    
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Active, &u.CreatedAt)
        users = append(users, u)
    }
    
    return users, rows.Err()
}

// IN 查询 / IN query
func getUsersByIDs(db *sql.DB, ids []int) ([]User, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    
    // 构建占位符 / Build placeholders
    placeholders := make([]string, len(ids))
    args := make([]interface{}, len(ids))
    for i, id := range ids {
        placeholders[i] = "?"
        args[i] = id
    }
    
    query := fmt.Sprintf(
        "SELECT id, name, email, age FROM users WHERE id IN (%s)",
        strings.Join(placeholders, ","),
    )
    
    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age)
        users = append(users, u)
    }
    
    return users, rows.Err()
}

// 带超时的查询 / Query with timeout
func queryWithTimeout(db *sql.DB) ([]User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    rows, err := db.QueryContext(ctx, "SELECT id, name FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name)
        users = append(users, u)
    }
    
    return users, rows.Err()
}

// 聚合查询 / Aggregate queries
func getStatistics(db *sql.DB) error {
    var count int
    var avgAge float64
    var maxAge, minAge int
    
    row := db.QueryRow(`
        SELECT COUNT(*), AVG(age), MAX(age), MIN(age) 
        FROM users WHERE active = true
    `)
    
    err := row.Scan(&count, &avgAge, &maxAge, &minAge)
    if err != nil {
        return err
    }
    
    fmt.Printf("Count: %d, Avg Age: %.2f, Max: %d, Min: %d\n",
        count, avgAge, maxAge, minAge)
    
    return nil
}

// 分页查询 / Pagination query
func getUsersPaginated(db *sql.DB, page, pageSize int) ([]User, int, error) {
    // 获取总数 / Get total count
    var total int
    db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
    
    // 查询当前页 / Query current page
    offset := (page - 1) * pageSize
    rows, err := db.Query(
        "SELECT id, name, email, age FROM users LIMIT ? OFFSET ?",
        pageSize, offset,
    )
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var u User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age)
        users = append(users, u)
    }
    
    return users, total, rows.Err()
}
```

## 7. 性能优化与最佳实践 / Performance & Best Practices

```go
package main

import (
    "database/sql"
    "log"
)

// 1. 使用预处理语句 / Use prepared statements
func optimizedBatchInsert(db *sql.DB, users []User) error {
    stmt, err := db.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, u := range users {
        _, err := stmt.Exec(u.Name, u.Email.String, u.Age)
        if err != nil {
            log.Printf("Failed to insert user %s: %v", u.Name, err)
        }
    }
    
    return nil
}

// 2. 正确关闭资源 / Properly close resources
func correctResourceHandling(db *sql.DB) {
    // 总是检查 rows.Close() 和 rows.Err()
    rows, err := db.Query("SELECT * FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()  // 确保关闭
    
    for rows.Next() {
        // 处理行...
    }
    
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}

// 3. 避免 SQL 注入 / Avoid SQL injection
func safeQuery(db *sql.DB, userInput string) {
    // 错误方式 - SQL 注入风险
    // db.Query("SELECT * FROM users WHERE name = '" + userInput + "'")
    
    // 正确方式 - 使用参数化查询
    db.Query("SELECT * FROM users WHERE name = ?", userInput)
}

// 4. 使用连接池配置 / Use connection pool configuration
func configurePool(db *sql.DB) {
    // 根据应用负载调整
    db.SetMaxOpenConns(25)      // 根据数据库服务器配置调整
    db.SetMaxIdleConns(5)       // 保持适量空闲连接
    db.SetConnMaxLifetime(300)  // 定期刷新连接
}

// 5. 批量操作使用事务 / Use transactions for batch operations
func batchWithTransaction(db *sql.DB, users []User) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare("INSERT INTO users (name, email, age) VALUES (?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, u := range users {
        _, err := stmt.Exec(u.Name, u.Email.String, u.Age)
        if err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```
