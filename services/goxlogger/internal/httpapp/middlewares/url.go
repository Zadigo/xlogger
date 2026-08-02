package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

const todoUrlParam = "todoId"

func TodoMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		serviceUuid := chi.URLParam(r, todoUrlParam)

		var ctx context.Context

		ctx = context.WithValue(r.Context(), todoUrlParam, serviceUuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
