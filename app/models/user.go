package models

import (
	"fmt"
)

type User struct {
	BaseModel
	UserId   uint64 `gorm:"primary_key;column:user_id"`
	Nickname string `gorm:"column:nickname"`
	Password string `gorm:"column:password"`
	Phone   string `gorm:"column:phone"`
	Gender   uint8  `gorm:"column:gender"`
	Platform uint8  `gorm:"column:platform"`
	Email    string `gorm:"column:email"`
	Avatar   string `gorm:"column:avatar"`
	TableId  uint64
}

// 分表
func (u User) TableName() string {
	return fmt.Sprintf("user_%d", getTableIndex(u.UserId))
}
func getTableIndex(uid uint64) int {
	// TODO：这儿自定义分表算法，这儿只是简单的对10取模
	return int(uid % 10)
}
