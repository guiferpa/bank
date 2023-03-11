package rest

import (
	"errors"
	"fmt"
	"github/guiferpa/bank/domain/account"
	"net/http"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func init() {
	httpin.UseGochiURLParam("path", chi.URLParam)
}

func httpRequestParamsErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	var invalidFieldError *httpin.InvalidFieldError

	render.Status(r, http.StatusBadRequest)

	if errors.As(err, &invalidFieldError) {
		fmt.Println(invalidFieldError)

		render.Respond(w, r, account.NewHandlerInvalidParamError(account.HandlerInvalidPathParam, "invalid path parameter", invalidFieldError.Field))
		return
	}

	render.Respond(w, r, account.NewHandlerError(account.HandlerUnknwonErrorCode, err.Error()))
}

func NewHTTPHandler(usecase account.UseCase) http.Handler {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))

	httpin.ReplaceDefaultErrorHandler(httpRequestParamsErrorHandler)

	router.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/accounts", func(r chi.Router) {
			r.Post("/", CreateAccount(usecase))
			r.With(httpin.NewInput(GetAccountByIDRequestParams{})).Get("/{id}", GetAccountByID(usecase))
		})
	})

	return router
}
