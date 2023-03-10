package account

type ErrorCode string

const (
	UseCaseDuplicatedAccountErrorCode ErrorCode = "domain.1"
	UseCaseUnknownErrorCode           ErrorCode = "domain.2"
)

type UseCaseCreateAccountError struct {
	Code    ErrorCode
	Message string
}

func (err *UseCaseCreateAccountError) Error() string {
	return err.Message
}

func NewUseCaseCreateAccountError(errorCode ErrorCode, message string) *UseCaseCreateAccountError {
	return &UseCaseCreateAccountError{errorCode, message}
}

const (
	StorageAccountNotFoundErrorCode ErrorCode = "infra.1"
)

type StorageRepositoryGetAccountByIDError struct {
	Code    ErrorCode
	Message string
}

func (err *StorageRepositoryGetAccountByIDError) Error() string {
	return err.Message
}

func NewStorageRepositoryGetAccountByIDError(errorCode ErrorCode, message string) *StorageRepositoryGetAccountByIDError {
	return &StorageRepositoryGetAccountByIDError{errorCode, message}
}
