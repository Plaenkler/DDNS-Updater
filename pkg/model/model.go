package model

import (
	"gorm.io/gorm"
)

type SyncJob struct {
	gorm.Model
	Provider    string    `gorm:"not null"`
	Domain      string    `gorm:"not null;unique"`
	User        string    `gorm:"not null"`
	Password    string    `gorm:"not null"`
	IPAddress   IPAddress `gorm:"foreignKey:IPAddressID"`
	IPAddressID *uint64
}

type IPAddress struct {
	gorm.Model
	Address string `gorm:"not null;unique"`
}

type Config struct {
	Port     uint64 `yaml:"Port"`
	Interval uint64 `yaml:"Interval"`
}
