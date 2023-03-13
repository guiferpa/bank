package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github/guiferpa/bank/domain/account"
	"github/guiferpa/bank/domain/log"

	"github.com/ggicci/httpin"
	"github.com/go-chi/render"
	"github.com/guiferpa/gody/v2"
	"github.com/guiferpa/gody/v2/rule"
)

type CreateAccountRequestBody struct {
	DocumentNumber string `json:"document_number" validate:"not_empty"`
}

type CreateAccountResponseBody struct {
	ID             uint   `json:"id"`
	DocumentNumber string `json:"document_number"`
}

func CreateAccount(usecase account.UseCase, logger log.LoggerRepository) http.HandlerFunc {
	validator := gody.NewValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CreateAccountRequestBody
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			render.Status(r, http.StatusBadRequest)

			if err == io.EOF {
				render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, "missing request body"))
				return
			}

			if _, ok := err.(*json.SyntaxError); ok {
				render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, "invalid request body"))
				return
			}

			if cerr, ok := err.(*json.UnmarshalTypeError); ok {
				render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, "wrong type", cerr.Field))
				return
			}

			render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, err.Error()))
			return
		}
		defer r.Body.Close()

		if err := validator.AddRules(rule.NotEmpty); err != nil {
			logger.Error(r.Context(), err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.Respond(w, r, account.NewHandlerError(account.HandlerUnknwonErrorCode, err.Error()))
			return
		}

		if _, err := validator.Validate(body); err != nil {
			render.Status(r, http.StatusUnprocessableEntity)

			if cerr, ok := err.(*rule.ErrNotEmpty); ok {
				render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, cerr.Error(), cerr.Field))
				return
			}

			render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, "", err.Error()))
			return
		}

		options := account.CreateAccountOptions{
			DocumentNumber: body.DocumentNumber,
		}
		accountID, err := usecase.CreateAccount(options)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)

			if cerr, ok := err.(*account.DomainError); ok && cerr.Code == account.DomainAccountAlreadyExistsErrorCode {
				render.Status(r, http.StatusConflict)
			}

			render.Respond(w, r, err)
			return
		}

		render.Status(r, http.StatusCreated)

		render.Respond(w, r, CreateAccountResponseBody{
			ID:             accountID,
			DocumentNumber: options.DocumentNumber,
		})

		logger.Info(r.Context(), "account created successful")
	})
}

type GetAccountByIDRequestParams struct {
	AccountID uint `in:"path=id"`
}

type GetAccountByIDResponseBody struct {
	ID             uint   `json:"id"`
	DocumentNumber string `json:"document_number"`
}

func GetAccountByID(usecase account.UseCase, logger log.LoggerRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(httpin.Input).(*GetAccountByIDRequestParams)

		acc, err := usecase.GetAccountByID(params.AccountID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)

			if cerr, ok := err.(*account.InfraError); ok && cerr.Code == account.InfraAccountNotFoundErrorCode {
				render.Status(r, http.StatusNotFound)
			}

			render.Respond(w, r, err)
			return
		}

		render.Status(r, http.StatusOK)

		render.Respond(w, r, GetAccountByIDResponseBody{
			ID:             acc.ID,
			DocumentNumber: acc.DocumentNumber,
		})

		logger.Info(r.Context(), "account retrieved by id successful")
	}
}

type CreateAccountTransactionRequestBody struct {
	AccountID       uint    `json:"account_id" validate:"min=0"`
	OperationTypeID uint    `json:"operation_type_id" validate:"min=0"`
	Amount          float64 `json:"amount" validate:"not_zero"`
}

type CreateAccountTransactionResponseBody struct {
	ID uint `json:"id"`
}

type NotZeroError struct {
	Field string
}

func (err *NotZeroError) Error() string {
	return "this value can't be zero"
}

type NotZeroRule struct{}

func (r *NotZeroRule) Name() string {
	return "not_zero"
}

func (r *NotZeroRule) Validate(field, value, _ string) (bool, error) {
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return true, &NotZeroError{field}
	}

	if amount == 0 {
		return true, &NotZeroError{field}
	}

	return true, nil
}

func CreateAccountTransaction(usecase account.UseCase, logger log.LoggerRepository) http.HandlerFunc {
	validator := gody.NewValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CreateAccountTransactionRequestBody
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			render.Status(r, http.StatusBadRequest)

			if err == io.EOF {
				render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, "missing request body"))
				return
			}

			if _, ok := err.(*json.SyntaxError); ok {
				render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, "invalid request body"))
				return
			}

			if cerr, ok := err.(*json.UnmarshalTypeError); ok {
				render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, "wrong type", cerr.Field))
				return
			}

			render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, err.Error()))
			return
		}
		defer r.Body.Close()

		if err := validator.AddRules(rule.Min, &NotZeroRule{}); err != nil {
			logger.Error(r.Context(), err.Error())
			render.Status(r, http.StatusInternalServerError)
			render.Respond(w, r, account.NewHandlerError(account.HandlerUnknwonErrorCode, err.Error()))
			return
		}

		if _, err := validator.Validate(body); err != nil {
			render.Status(r, http.StatusUnprocessableEntity)

			if cerr, ok := err.(*rule.ErrMin); ok {
				render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, cerr.Error(), cerr.Field))
				return
			}

			if cerr, ok := err.(*NotZeroError); ok {
				render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, cerr.Error(), cerr.Field))
				return
			}

			render.Respond(w, r, account.NewHandlerInvalidFieldError(account.HandlerInvalidPayloadErrorCode, "", err.Error()))
			return
		}

		options := account.CreateTransactionOptions{
			AccountID:       body.AccountID,
			OperationTypeID: body.OperationTypeID,
			Amount:          int64(body.Amount * 100),
			EventDate:       time.Now(),
		}
		transID, err := usecase.CreateTransaction(options)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)

			if cerr, ok := err.(*account.DomainError); ok {
				if cerr.Code == account.DomainOperationTypeDoesntExistErrorCode {
					render.Status(r, http.StatusNotFound)
				}
			}

			if cerr, ok := err.(*account.InfraError); ok && cerr.Code == account.InfraAccountNotFoundErrorCode {
				render.Status(r, http.StatusNotFound)
			}

			render.Respond(w, r, err)
			return
		}

		render.Status(r, http.StatusCreated)

		render.Respond(w, r, CreateAccountTransactionResponseBody{
			ID: transID,
		})

		logger.Info(r.Context(), "account transaction created successful")
	})
}
