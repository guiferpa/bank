package account

type ErrorCode string

const (
	HandlerUnknwonErrorCode        ErrorCode = "handler.1"
	HandlerInvalidPayloadErrorCode ErrorCode = "handler.2"
	HandlerBadRequestErrorCode     ErrorCode = "handler.3"
	HandlerInvalidPathParam        ErrorCode = "handler.4"
)

type HandlerError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewHandlerError(errorCode ErrorCode, message string) *HandlerError {
	return &HandlerError{errorCode, message}
}

type HandlerInvalidFieldError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
}

func (err *HandlerInvalidFieldError) Error() string {
	return err.Message
}

func NewHandlerInvalidFieldError(errorCode ErrorCode, message, field string) *HandlerInvalidFieldError {
	return &HandlerInvalidFieldError{errorCode, message, field}
}

type HandlerInvalidParamError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Param   string    `json:"parameter,omitempty"`
}

func NewHandlerInvalidParamError(errorCode ErrorCode, message, param string) *HandlerInvalidParamError {
	return &HandlerInvalidParamError{errorCode, message, param}
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
