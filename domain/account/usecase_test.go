package account

import "testing"

type MockStorageRepository struct {
	NCalledCreateAccount              int
	NCalledCreatedTransaction         int
	NCalledGetAccountByID             int
	NCalledHasAccountByDocumentNumber int
	DocumentNumberResult              string
	HasAccountByDocumentNumberResult  bool
	GetAccountByIDErrorResult         error
}

func (msr *MockStorageRepository) CreateAccount(opts CreateAccountOptions) (uint, error) {
	msr.NCalledCreateAccount += 1
	msr.DocumentNumberResult = opts.DocumentNumber
	return 0, nil
}

func (msr *MockStorageRepository) CreateTransaction(opts CreateTransactionOptions) (uint, error) {
	msr.NCalledCreatedTransaction += 1
	return 0, nil
}

func (msr *MockStorageRepository) GetAccountByID(accountID uint) (Account, error) {
	msr.NCalledGetAccountByID += 1
	return Account{}, msr.GetAccountByIDErrorResult
}

func (msr *MockStorageRepository) HasAccountByDocumentNumber(documentNumber string) (bool, error) {
	msr.NCalledHasAccountByDocumentNumber += 1
	return msr.HasAccountByDocumentNumberResult, nil
}

func TestCreateAccount(t *testing.T) {
	suite := []struct {
		DocumentNumber                            string
		ExpectedNCalledHasAccountByDocumentNumber int
		ExpectedNCalledCreateAccount              int
		ExpectedDocumentNumberResult              string
	}{
		{
			DocumentNumber: "123",
			ExpectedNCalledHasAccountByDocumentNumber: 1,
			ExpectedNCalledCreateAccount:              1,
			ExpectedDocumentNumberResult:              "123",
		},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{}
		svc := &UseCaseService{storage: mock}

		opts := CreateAccountOptions{DocumentNumber: s.DocumentNumber}
		if _, err := svc.CreateAccount(opts); err != nil {
			t.Error(err)
			return
		}

		if got, expected := mock.NCalledHasAccountByDocumentNumber, s.ExpectedNCalledHasAccountByDocumentNumber; got != expected {
			t.Errorf("unexpected N called HasAccountByDocumentNumber, got: %v, expected: %v", got, expected)
			return
		}

		if got, expected := mock.NCalledCreateAccount, s.ExpectedNCalledCreateAccount; got != expected {
			t.Errorf("unexpected N called CreateAccount, got: %v, expected: %v", got, expected)
			return
		}

		if got, expected := mock.DocumentNumberResult, s.ExpectedDocumentNumberResult; got != expected {
			t.Errorf("unexpected document number, got: %v, expected: %v", got, expected)
			return
		}
	}
}

func TestCreateAccountWithDocumentNumberAlreadyRegistered(t *testing.T) {
	suite := []struct {
		DocumentNumber                            string
		ExpectedNCalledHasAccountByDocumentNumber int
		ExpectedNCalledCreateAccount              int
		HasAccountByDocumentNumberResult          bool
	}{
		{
			DocumentNumber: "123",
			ExpectedNCalledHasAccountByDocumentNumber: 1,
			ExpectedNCalledCreateAccount:              0,
			HasAccountByDocumentNumberResult:          true,
		},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{
			HasAccountByDocumentNumberResult: s.HasAccountByDocumentNumberResult,
		}
		svc := &UseCaseService{storage: mock}

		opts := CreateAccountOptions{DocumentNumber: s.DocumentNumber}
		_, err := svc.CreateAccount(opts)

		cerr, ok := err.(*UseCaseCreateAccountError)

		if !ok {
			t.Error("unexpected error")
			return
		}

		if got, expected := cerr.Code, UseCaseCreateAccountDuplicatedAccountErrorCode; got != expected {
			t.Errorf("unexpected error code, got: %v, expected: %v", got, expected)
			return
		}

		if got, expected := mock.NCalledHasAccountByDocumentNumber, s.ExpectedNCalledHasAccountByDocumentNumber; got != expected {
			t.Errorf("unexpected N called HasAccountByDocumentNumber, got: %v, expected: %v", got, expected)
			return
		}

		if got, expected := mock.NCalledCreateAccount, s.ExpectedNCalledCreateAccount; got != expected {
			t.Errorf("unexpected N called CreateAccount, got: %v, expected: %v", got, expected)
			return
		}
	}
}

func TestCreateTransaction(t *testing.T) {
	suite := []struct {
		ExpectedNCalledGetAccountByID    int
		ExpectedNCalledCreateTransaction int
	}{
		{
			ExpectedNCalledGetAccountByID:    1,
			ExpectedNCalledCreateTransaction: 1,
		},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{}
		svc := &UseCaseService{storage: mock}

		opts := CreateTransactionOptions{}
		if _, err := svc.CreateTransaction(opts); err != nil {
			t.Error(err)
			return
		}

		if got, expected := mock.NCalledGetAccountByID, s.ExpectedNCalledGetAccountByID; got != expected {
			t.Errorf("unexpected N called GetAccountByID, got: %v, expected: %v", got, expected)
			return
		}

		if got, expected := mock.NCalledCreatedTransaction, s.ExpectedNCalledCreateTransaction; got != expected {
			t.Errorf("unexpected N called CreateTransaction, got: %v, expected: %v", got, expected)
			return
		}
	}
}

func TestGetAccountById(t *testing.T) {
	suite := []struct {
		ExpectedNCalledGetAccountByID int
	}{
		{ExpectedNCalledGetAccountByID: 1},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{}
		svc := &UseCaseService{storage: mock}

		accountID := uint(20)
		if _, err := svc.GetAccountByID(accountID); err != nil {
			t.Error(err)
			return
		}

		if got, expected := mock.NCalledGetAccountByID, s.ExpectedNCalledGetAccountByID; got != expected {
			t.Errorf("unexpected N called GetAccountByID, got: %v, expected: %v", got, expected)
			return
		}
	}
}

func TestGetAccountByIdWithNotFound(t *testing.T) {
	suite := []struct {
		ExpectedNCalledGetAccountByID int
		GetAccountByIDErrorResult     error
	}{
		{
			ExpectedNCalledGetAccountByID: 1,
			GetAccountByIDErrorResult:     NewStorageRepositoryGetAccountByIDError(StorageAccountNotFoundErrorCode, "account not found"),
		},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{
			GetAccountByIDErrorResult: s.GetAccountByIDErrorResult,
		}
		svc := &UseCaseService{storage: mock}

		_, err := svc.GetAccountByID(20)
		cerr, ok := err.(*StorageRepositoryGetAccountByIDError)
		if !ok {
			t.Error("unexpected error")
			return
		}

		if got, expected := cerr.Code, StorageAccountNotFoundErrorCode; got != expected {
			t.Errorf("unexpected error code, got: %v, expected: %v", got, expected)
			return
		}

		if got, expected := mock.NCalledGetAccountByID, s.ExpectedNCalledGetAccountByID; got != expected {
			t.Errorf("unexpected N called GetAccountByID, got: %v, expected: %v", got, expected)
			return
		}
	}
}
