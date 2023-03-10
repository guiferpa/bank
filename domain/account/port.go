package account

import "time"

type CreateAccountOptions struct {
	DocumentNumber string
}

type CreateTransactionOptions struct {
	AccountID       uint
	OperationTypeID uint
	Amount          int64
	EventDate       time.Time
}

type StorageRepository interface {
	CreateAccount(CreateAccountOptions) (uint, error)
	GetAccountByID(uint) (Account, error)
	HasAccountByDocumentNumber(string) (bool, error)
	CreateTransaction(CreateTransactionOptions) (uint, error)
}

type UseCase interface {
	CreateAccount(CreateAccountOptions) (uint, error)
	GetAccountByID(uint) (Account, error)
	CreateTransaction(CreateTransactionOptions) (uint, error)
}
