package model

import (
	"time"
)

// User ORM 模型
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(64);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	Phone     string    `gorm:"type:varchar(20);index" json:"phone"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
