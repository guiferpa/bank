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

type StorageRepositoryCreateAccountError struct {
	Message string
}

func (err *StorageRepositoryCreateAccountError) Error() string {
	return err.Message
}

func NewStorageRepositoryCreateAccountError(message string) *StorageRepositoryCreateAccountError {
	return &StorageRepositoryCreateAccountError{message}
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
