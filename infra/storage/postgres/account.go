package postgres

import "time"

type Account struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	DocumentNumber string    `gorm:"index;unique"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (a *Account) TableName() string {
	return "accounts"
}