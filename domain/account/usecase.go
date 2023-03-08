package account

type UseCaseService struct {
	storage StorageRepository
}

func (ucs *UseCaseService) CreateAccount(opts CreateAccountOptions) (int, error) {
	return ucs.storage.CreateAccount(opts)
}
