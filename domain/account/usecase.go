package account

type UseCaseService struct {
	storage StorageRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (uint, error) {
	id, err := ucs.storage.CreateAccount(opts)
	if err != nil {
		return 0, NewUseCaseCreateAccountError(err.Error())
	}

	return id, nil
}

func (ucs *UseCaseService) CreateTransaction(opts CreateTransactionOptions) (uint, error) {
	return ucs.storage.CreateTransaction(opts)
}

func (ucs *UseCaseService) GetAccountByID(accountID uint) (Account, error) {
	return ucs.storage.GetAccountByID(accountID)
}
