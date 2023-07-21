package model

import (
	"gorm.io/gorm"
)

type SyncJob struct {
	gorm.Model
	Provider    string    `gorm:"not null"`
	Params      string    `gorm:"not null"`
	IPAddress   IPAddress `gorm:"foreignKey:IPAddressID"`
	IPAddressID *uint64
}

type IPAddress struct {
	gorm.Model
	Address string `gorm:"not null;unique"`
}

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	EMail    string `gorm:"not null;unique"`
}
