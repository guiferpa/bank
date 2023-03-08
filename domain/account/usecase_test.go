package account

import "testing"

type MockStorageRepository struct {
	NCalledCreateAccount      int
	NCalledCreatedTransaction int
	DocumentNumberResult      string
}

func (msr *MockStorageRepository) CreateAccount(opts CreateAccountOptions) (int, error) {
	msr.NCalledCreateAccount += 1
	msr.DocumentNumberResult = opts.DocumentNumber
	return 0, nil
}

func (msr *MockStorageRepository) CreateTransaction(opts CreateTransactionOptions) (int, error) {
	msr.NCalledCreatedTransaction += 1
	return 0, nil
}

func TestCreateAccount(t *testing.T) {
	suite := []struct {
		DocumentNumber               string
		ExpectedNCalledCreateAccount int
		ExpectedDocumentNumberResult string
	}{
		{DocumentNumber: "123", ExpectedNCalledCreateAccount: 1, ExpectedDocumentNumberResult: "123"},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{}
		svc := &UseCaseService{storage: mock}

		opts := CreateAccountOptions{DocumentNumber: s.DocumentNumber}
		if _, err := svc.CreateAccount(opts); err != nil {
			t.Error(err)
			return
		}

		if got, expected := mock.DocumentNumberResult, s.ExpectedDocumentNumberResult; got != expected {
			t.Errorf("unexpected document number, got: %v, expected: %v", got, expected)
			return
		}
	}
}

func TestCreateTransaction(t *testing.T) {
	suite := []struct {
		ExpectedNCalledCreateTransaction int
	}{
		{ExpectedNCalledCreateTransaction: 1},
	}

	for _, s := range suite {
		mock := &MockStorageRepository{}
		svc := &UseCaseService{storage: mock}

		opts := CreateTransactionOptions{}
		if _, err := svc.CreateTransaction(opts); err != nil {
			t.Error(err)
			return
		}

		if got, expected := mock.NCalledCreatedTransaction, s.ExpectedNCalledCreateTransaction; got != expected {
			t.Errorf("unexpected N called CreateTransaction, got: %v, expected: %v", got, expected)
			return
		}
	}
}
