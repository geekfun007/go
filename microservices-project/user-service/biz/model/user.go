package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 不序列化密码
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Age       int            `gorm:"type:int;default:0" json:"age"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 软删除
}

func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM 钩子：创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里添加创建前的逻辑
	return nil
}

// BeforeUpdate GORM 钩子：更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里添加更新前的逻辑
	return nil
}
