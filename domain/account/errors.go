package account

type ErrorCode string

const (
	HandlerUnknwonErrorCode        ErrorCode = "handler.1"
	HandlerInvalidPayloadErrorCode ErrorCode = "handler.2"
	HandlerBadRequestErrorCode     ErrorCode = "handler.3"
)

type HandlerError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewHandlerError(errorCode ErrorCode, message string) *HandlerError {
	return &HandlerError{errorCode, message}
}

type HandlerInvalidPayloadError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
}

func NewHandlerInvalidPayloadError(errorCode ErrorCode, message, field string) *HandlerInvalidPayloadError {
	return &HandlerInvalidPayloadError{errorCode, message, field}
}

const (
	DomainAccountAlreadyExistsErrorCode ErrorCode = "domain.1"
)

type DomainError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (err *DomainError) Error() string {
	return err.Message
}

func NewDomainError(errorCode ErrorCode, message string) *DomainError {
	return &DomainError{errorCode, message}
}

const (
	InfraUnknownError             ErrorCode = "infra.1"
	InfraAccountNotFoundErrorCode ErrorCode = "infra.2"
)

type InfraError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (err *InfraError) Error() string {
	return err.Message
}

func NewInfraError(errorCode ErrorCode, message string) *InfraError {
	return &InfraError{errorCode, message}
}
