package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github/guiferpa/bank/domain/account"

	"github.com/go-chi/render"
	"github.com/guiferpa/gody/v2"
	"github.com/guiferpa/gody/v2/rule"
)

type CreateAccountHTTPRequest struct {
	DocumentNumber string `json:"document_number" validate:"not_empty"`
}

type CreateAccountHTTPResponse struct {
	ID             uint   `json:"id"`
	DocumentNumber string `json:"document_number"`
}

func CreateAccount(usecase account.UseCase) http.HandlerFunc {
	validator := gody.NewValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CreateAccountHTTPRequest
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			if err == io.EOF {
				render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, "missing request body"))
				return
			}

			if _, ok := err.(*json.SyntaxError); ok {
				render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, "invalid request body"))
				return
			}

			if cerr, ok := err.(*json.UnmarshalTypeError); ok {
				render.Respond(w, r, account.NewHandlerInvalidPayloadError(account.HandlerInvalidPayloadErrorCode, "wrong type", cerr.Field))
				return
			}

			render.Respond(w, r, account.NewHandlerError(account.HandlerBadRequestErrorCode, err.Error()))
			return
		}
		defer r.Body.Close()

		if _, err := validator.Validate(body); err != nil {
			if cerr, ok := err.(*rule.ErrNotEmpty); ok {
				render.Respond(w, r, account.NewHandlerInvalidPayloadError(account.HandlerInvalidPayloadErrorCode, cerr.Error(), cerr.Field))
				return
			}

			render.Respond(w, r, account.NewHandlerInvalidPayloadError(account.HandlerInvalidPayloadErrorCode, "", err.Error()))
			return
		}

		options := account.CreateAccountOptions{
			DocumentNumber: body.DocumentNumber,
		}
		accountID, err := usecase.CreateAccount(options)
		if err != nil {
			render.Respond(w, r, err)
			return
		}

		render.Respond(w, r, CreateAccountHTTPResponse{
			ID:             accountID,
			DocumentNumber: options.DocumentNumber,
		})
	})
}
