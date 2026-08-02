package httpapp

import (
	"time"

	"github.com/Zadigo/goxlogger/internal/handlers"
	"github.com/Zadigo/goxlogger/internal/httpapp/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	a.router = chi.NewRouter()

	a.router.Use(middleware.RequestID)
	a.router.Use(middleware.RealIP)
	a.router.Use(middlewares.Cors)
	a.router.Use(middlewares.Authorization)
	a.router.Use(middleware.AllowContentType("application/json"))
	a.router.Use(middleware.Throttle(1000))
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middlewares.JsonHeartbeat("/health"))
	a.router.Use(middleware.Timeout(60 * time.Second))

	a.router.Route("/v1/files", a.loadBaseRoutes)
}

func (a *App) loadBaseRoutes(router chi.Router) {
	baseHandlers := &handlers.BaseRouteHandlers{}

	baseHandlers.SetApp(a)

	router.Use(middlewares.TodoMiddleware)

	router.Get("/", baseHandlers.GetFiles)
	router.Get("/files/${fileId}", baseHandlers.GetLogs)
}
