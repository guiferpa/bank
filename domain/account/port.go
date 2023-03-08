package account

import "time"

type CreateAccountOptions struct {
	DocumentNumber string
}

type CreateTransactionOptions struct {
	AccountID       int
	OperationTypeID int
	Amount          int
	EventDate       time.Time
}

type StorageRepository interface {
	CreateAccount(CreateAccountOptions) (int, error)
	CreateTransaction(CreateTransactionOptions) (int, error)
	GetAccountByID(int) (Account, error)
}

type UseCase interface {
	CreateAccount(CreateAccountOptions) (int, error)
	CreateTransaction(CreateTransactionOptions) (int, error)
	GetAccountByID(int) (Account, error)
}
