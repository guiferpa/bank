package postgres

import (
	"gorm.io/gorm"
)

type OperationType struct {
	gorm.Model

	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Description string `gorm:"size:128"`
}

func (ot *OperationType) TableName() string {
	return "operation_types"
}

var OperationTypeSeedData = []OperationType{
	{ID: 1, Description: "COMPRA A VISTA"},
	{ID: 2, Description: "COMPRA PARCELADA"},
	{ID: 3, Description: "SAQUE"},
	{ID: 4, Description: "PAGAMENTO"},
}
