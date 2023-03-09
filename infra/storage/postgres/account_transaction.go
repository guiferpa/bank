package postgres

import (
	"time"

	"gorm.io/gorm"
)

type AccountTransaction struct {
	gorm.Model

	ID              uint `gorm:"primaryKey;autoIncrement"`
	AccountID       uint
	OperationTypeID uint
	Amount          int64
	EventDate       time.Time

	Account       Account       `gorm:"foreignKey:AccountID"`
	OperationType OperationType `gorm:"foreignKey:OperationTypeID"`
}

func (at *AccountTransaction) TableName() string {
	return "transactions"
}
