package models

import "time"

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//    type User struct {
//      BaseModel
//    }
type BaseModel struct {
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	DeletedAt *time.Time `sql:"index"`
}
