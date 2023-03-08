package account

type UseCaseService struct {
	storage StorageRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (uint, error) {
	return ucs.storage.CreateAccount(opts)
}

func (ucs *UseCaseService) CreateTransaction(opts CreateTransactionOptions) (uint, error) {
	return ucs.storage.CreateTransaction(opts)
}

func (ucs *UseCaseService) GetAccountByID(accountID uint) (Account, error) {
	return ucs.storage.GetAccountByID(accountID)
}
