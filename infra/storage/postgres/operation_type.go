package postgres

import "time"

type OperationType struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Description string    `gorm:"size:128"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
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
