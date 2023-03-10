package rest

import (
	"github/guiferpa/bank/domain/account"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewHTTPHandler(usecase account.UseCase) http.Handler {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/accounts", func(r chi.Router) {
			r.Post("/", CreateAccount(usecase))
		})
	})

	return router
}
