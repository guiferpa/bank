package account

type UseCaseService struct {
	storage StorageRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (int, error) {
	return ucs.storage.CreateAccount(opts)
}

func (ucs *UseCaseService) CreateTransaction(opts CreateTransactionOptions) (int, error) {
	return ucs.storage.CreateTransaction(opts)
}
