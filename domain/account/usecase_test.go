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

		if _, err := svc.CreateAccount(CreateAccountOptions{DocumentNumber: s.DocumentNumber}); err != nil {
			t.Error(err)
			return
		}

		if got, expected := mock.DocumentNumberResult, s.ExpectedDocumentNumberResult; got != expected {
			t.Errorf("unexpected document number, got: %v, expected: %v", got, expected)
			return
		}
	}
}
