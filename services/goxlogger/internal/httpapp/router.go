package httpapp

import (
	"time"

	"github.com/Zadigo/goxlogger/internal/handlers"
	"github.com/Zadigo/goxlogger/internal/httpapp/middlewares"
	"github.com/Zadigo/goxlogger/internal/mcpapp"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	a.router = chi.NewRouter()

	a.router.Use(middlewares.SecurityHeaders)
	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middlewares.Cors)
	a.router.Use(middlewares.Authorization)
	a.router.Use(middlewares.RateLimitLogin)
	a.router.Use(middleware.AllowContentType("application/json"))
	a.router.Use(middleware.Throttle(1000))
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middlewares.JsonHeartbeat("/health"))
	a.router.Use(middleware.Timeout(60 * time.Second))

	a.router.Route("/v1/", a.loadBaseRoutes)
}

func (a *App) loadBaseRoutes(r chi.Router) {
	baseHandlers := &handlers.BaseRouteHandlers{}
	baseHandlers.SetApp(a)

	mcpApp := &mcpapp.McpApp{}
	r.Route("/mcp", func(r chi.Router) {
		r.Handle("/*", mcpApp.ChiHandle())
	})

	// r.Get("/files", baseHandlers.GetFilesHandler)

	// r.Route("/{fileId}", func(r chi.Router) {
	// 	r.Use(middlewares.UrlParamMiddleware("fileId"))
	// 	r.Get("/", baseHandlers.GetLogsHandler)
	// })
}
