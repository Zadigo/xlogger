package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

const fileIdUrlParam = "fileId"

func FileIdMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fileId := chi.URLParam(r, fileIdUrlParam)

		var ctx context.Context

		ctx = context.WithValue(r.Context(), fileIdUrlParam, fileId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
