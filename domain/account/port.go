package account

import "time"

type CreateAccountOptions struct {
	DocumentNumber string
}

type CreateTransactionOptions struct {
	AccountID       uint
	OperationTypeID uint
	Amount          uint
	EventDate       time.Time
}

type StorageRepository interface {
	CreateAccount(CreateAccountOptions) (uint, error)
	CreateTransaction(CreateTransactionOptions) (uint, error)
	GetAccountByID(uint) (Account, error)
}

type UseCase interface {
	CreateAccount(CreateAccountOptions) (uint, error)
	CreateTransaction(CreateTransactionOptions) (uint, error)
	GetAccountByID(uint) (Account, error)
}
