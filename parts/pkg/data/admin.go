package data

import (
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID        uint64 `gorm:"primarykey" json:"id"`
	IP        string `json:"ip"`
	SignIn    bool
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
