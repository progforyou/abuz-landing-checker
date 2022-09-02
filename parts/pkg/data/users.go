package data

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	ID            uint64         `gorm:"primarykey" json:"id"`
	GeneratedHash string         `gorm:"size:50" json:"generated_hash"`
	IP            string         `json:"ip_address"`
	IPLocation    string         `json:"ip_location"`
	FP            string         `gorm:"size:256" json:"fp"`
	TelegramName  string         `gorm:"size:50" json:"telegram"`
	Count         uint64         `json:"count"`
	IsAntik       bool           `json:"is_antik"`
	UniqHash      string         `gorm:"size:100" json:"uniq_hash"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
