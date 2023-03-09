package postgres

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	ID             uint   `gorm:"primaryKey;autoIncrement"`
	DocumentNumber string `gorm:"index;unique"`
}

func (a *Account) TableName() string {
	return "accounts"
}
