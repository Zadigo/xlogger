package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zadigo/goxlogger/internal/handlers"
	"github.com/Zadigo/goxlogger/internal/httpapp"
)

func GenericRecorder(t *testing.T, method string, path string, handle func(w http.ResponseWriter, r *http.Request)) *httptest.ResponseRecorder {
	t.Helper()

	handler := http.HandlerFunc(handle)

	request := httptest.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}

func CreateGetFilesRecorder(t *testing.T) *httptest.ResponseRecorder {
	ctx := context.WithValue(t.Context(), "rootDir", "../")

	app := httpapp.NewApp(ctx)

	handlers := handlers.BaseRouteHandlers{}
	handlers.SetApp(app)

	return GenericRecorder(t, "GET", "/v1/files", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetFiles(w, r)
	})
}

func CreateGetLogsRecorder(t *testing.T) *httptest.ResponseRecorder {
	ctx := context.WithValue(t.Context(), "rootDir", "../")
	app := httpapp.NewApp(ctx)

	handlers := handlers.BaseRouteHandlers{}
	handlers.SetApp(app)

	// 1. Create a fresh ServeMux for the test
	mux := http.NewServeMux()

	// 2. Register the route pattern matching your production routing
	mux.HandleFunc("GET /v1/files/{filename}", handlers.GetLogs)

	// 3. Create the recorder and the request pointing to the specific asset
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/files/ZXhhbXBsZTEubG9n", nil)

	requestCtx := context.WithValue(req.Context(), "fileId", "ZXhhbXBsZTEubG9n")
	req = req.WithContext(requestCtx)

	// 4. Route the request through the mux so parameters are parsed
	mux.ServeHTTP(w, req)

	return w
}
