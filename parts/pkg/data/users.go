package data

import (
	"gorm.io/gorm"
	"time"
)

type Users struct {
	ID         uint64         `gorm:"primarykey" json:"id"`
	IP         string         `json:"ip"`
	IPLocation string         `json:"ip-location"`
	FP         string         `gorm:"size:256" json:"fp"`
	Count      uint64         `json:"count"`
	IsAntik    bool           `json:"is-antik"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
