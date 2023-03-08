package account

type UseCaseCreateAccountError struct {
	Message string
}

func (err *UseCaseCreateAccountError) Error() string {
	return err.Message
}

func NewUseCaseCreateAccountError(message string) *UseCaseCreateAccountError {
	return &UseCaseCreateAccountError{message}
}

type StorageRepositoryGetAccountByIDError struct {
	Message string
}

func (err *StorageRepositoryGetAccountByIDError) Error() string {
	return err.Message
}

func NewStorageRepositoryGetAccountByIDError(message string) *StorageRepositoryGetAccountByIDError {
	return &StorageRepositoryGetAccountByIDError{message}
}
