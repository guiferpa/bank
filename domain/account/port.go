package account

type CreateAccountOptions struct {
	DocumentNumber string
}

type StorageRepository interface {
	CreateAccount(CreateAccountOptions) (int, error)
}

type UseCase interface {
	CreateAccount(CreateAccountOptions) (int, error)
}
