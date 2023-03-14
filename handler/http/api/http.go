package api

import (
	"context"
	"errors"
	"fmt"
	"github/guiferpa/bank/domain/account"
	"github/guiferpa/bank/domain/log"
	"net/http"
	"time"

	"github.com/ggicci/httpin"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func init() {
	httpin.UseGochiURLParam("path", chi.URLParam)
}

func HTTPResponseLoggerMiddleware(logger log.LoggerRepository) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t := time.Now()

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			defer func() {
				elapsed := time.Since(t)
				ctxvalue, _ := r.Context().Value(log.LoggerContextKey).(*log.LoggerContext)
				value := &log.LoggerContext{
					RequestID: ctxvalue.RequestID,
					Payload: map[string]interface{}{
						"response_time": elapsed,
						"method":        r.Method,
						"byte_written":  ww.BytesWritten(),
						"status_code":   ww.Status(),
						"from_ip":       r.RemoteAddr,
						"proto":         r.Proto,
						"request_uri":   r.RequestURI,
						"host":          r.Host,
					},
				}
				logger.Info(context.WithValue(r.Context(), log.LoggerContextKey, value), fmt.Sprintf("%s - %d - %s://%s/%s in %s", r.Method, ww.Status(), scheme, r.Host, r.RequestURI, elapsed))
			}()

			h.ServeHTTP(ww, r)
		})
	}
}

func SetRequestContextMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = fmt.Sprintf("%v", time.Now().Unix())
		}

		value := &log.LoggerContext{RequestID: requestId}
		ctx := context.WithValue(r.Context(), log.LoggerContextKey, value)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewHTTPHandler(usecase account.UseCase, logger log.LoggerRepository) http.Handler {
	router := chi.NewRouter()

	router.Use(render.SetContentType(render.ContentTypeJSON), SetRequestContextMiddleware, HTTPResponseLoggerMiddleware(logger))

	httpin.ReplaceDefaultErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		var invalidFieldError *httpin.InvalidFieldError

		render.Status(r, http.StatusBadRequest)

		if errors.As(err, &invalidFieldError) {
			render.Respond(w, r, account.NewHandlerInvalidParamError(account.HandlerInvalidPathParam, "invalid path parameter", invalidFieldError.Field))
			return
		}

		render.Respond(w, r, account.NewHandlerError(account.HandlerUnknwonErrorCode, err.Error()))
	})

	router.Route("/api/v1", func(v1 chi.Router) {
		v1.Route("/accounts", func(r chi.Router) {
			r.Post("/", CreateAccount(usecase, logger))
			r.With(httpin.NewInput(GetAccountByIDRequestParams{})).Get("/{id}", GetAccountByID(usecase, logger))
			r.Post("/transaction", CreateAccountTransaction(usecase, logger))
		})
	})

	return router
}
