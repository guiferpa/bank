package account

type ErrorCode string

const (
	UseCaseCreateAccountDuplicatedAccountErrorCode ErrorCode = "domain.1"
	UseCaseCreateAccountUnknownErrorCode           ErrorCode = "domain.2"

	UseCaseGetAccountByIDUnknownErrorCode ErrorCode = "domain.3"
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

type UseCaseGetAccountByIDError struct {
	Code    ErrorCode
	Message string
}

func (err *UseCaseGetAccountByIDError) Error() string {
	return err.Message
}

func NewUseCaseGetAccountByIDError(errorCode ErrorCode, message string) *UseCaseGetAccountByIDError {
	return &UseCaseGetAccountByIDError{errorCode, message}
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
