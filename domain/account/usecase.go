package account

import (
	"github/guiferpa/bank/domain/log"
)

type UseCaseService struct {
	storage StorageRepository
	logger  log.LoggerRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (uint, error) {
	has, err := ucs.storage.HasAccountByDocumentNumber(opts.DocumentNumber)
	if err != nil {
		return 0, err
	}

	if has {
		return 0, NewDomainError(DomainAccountAlreadyExistsErrorCode, "account already exists")
	}

	accountID, err := ucs.storage.CreateAccount(opts)
	if err != nil {
		return 0, err
	}

	return accountID, nil
}

func (ucs *UseCaseService) CreateTransaction(opts CreateTransactionOptions) (uint, error) {
	if _, err := ucs.storage.GetAccountByID(opts.AccountID); err != nil {
		return 0, err
	}

	hasOperationType, err := ucs.storage.HasOperationTypeByID(opts.OperationTypeID)
	if err != nil {
		return 0, err
	}

	if !hasOperationType {
		return 0, NewDomainError(DomainOperationTypeDoesntExistErrorCode, "operation type doesn't exist")
	}

	transID, err := ucs.storage.CreateTransaction(opts)
	if err != nil {
		return 0, err
	}

	return transID, nil
}

func (ucs *UseCaseService) GetAccountByID(accountID uint) (Account, error) {
	acc, err := ucs.storage.GetAccountByID(accountID)
	if err != nil {
		return Account{}, err
	}

	return acc, nil
}

func NewUseCaseService(storage StorageRepository, logger log.LoggerRepository) *UseCaseService {
	return &UseCaseService{storage, logger}
}
