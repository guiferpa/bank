package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github/guiferpa/bank/domain/account"

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

func CreateAccount(usecase account.UseCase) http.HandlerFunc {
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
	})
}

type GetAccountByIDRequestParams struct {
	AccountID uint `in:"path=id"`
}

type GetAccountByIDResponseBody struct {
	ID             uint   `json:"id"`
	DocumentNumber string `json:"document_number"`
}

func GetAccountByID(usecase account.UseCase) http.HandlerFunc {
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
	}
}
