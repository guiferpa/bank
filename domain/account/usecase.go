package account

type UseCaseService struct {
	storage StorageRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (uint, error) {
	has, err := ucs.storage.HasAccountByDocumentNumber(opts.DocumentNumber)
	if err != nil {
		return 0, NewUseCaseCreateAccountError(UseCaseCreateAccountUnknownErrorCode, err.Error())
	}

	if has {
		return 0, NewUseCaseCreateAccountError(UseCaseCreateAccountDuplicatedAccountErrorCode, "account already exists")
	}

	accountID, err := ucs.storage.CreateAccount(opts)
	if err != nil {
		return 0, NewUseCaseCreateAccountError(UseCaseCreateAccountUnknownErrorCode, err.Error())
	}

	return accountID, nil
}

func (ucs *UseCaseService) CreateTransaction(opts CreateTransactionOptions) (uint, error) {
	return ucs.storage.CreateTransaction(opts)
}

func (ucs *UseCaseService) GetAccountByID(accountID uint) (Account, error) {
	acc, err := ucs.storage.GetAccountByID(accountID)
	if err != nil {
		if _, ok := err.(*StorageRepositoryGetAccountByIDError); ok {
			return Account{}, err
		}

		return Account{}, NewUseCaseGetAccountByIDError(UseCaseGetAccountByIDUnknownErrorCode, err.Error())
	}

	return acc, nil
}
