package account

type UseCaseService struct {
	storage StorageRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (uint, error) {
	has, err := ucs.storage.HasAccountByDocumentNumber(opts.DocumentNumber)
	if err != nil {
		return 0, NewUseCaseCreateAccountError(UseCaseUnknownErrorCode, err.Error())
	}

	if has {
		return 0, NewUseCaseCreateAccountError(UseCaseDuplicatedAccountErrorCode, "account already exists")
	}

	return ucs.storage.CreateAccount(opts)
}

func (ucs *UseCaseService) CreateTransaction(opts CreateTransactionOptions) (uint, error) {
	return ucs.storage.CreateTransaction(opts)
}

func (ucs *UseCaseService) GetAccountByID(accountID uint) (Account, error) {
	return ucs.storage.GetAccountByID(accountID)
}
