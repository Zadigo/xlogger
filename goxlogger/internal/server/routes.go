package server

import (
	"time"

	"github.com/Zadigo/goxlogger/internal/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Middlewares()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(handlers.Cors)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/logs", a.loadBaseRoutes)

	a.router = router
}

func (a *App) loadBaseRoutes(r chi.Router) {
	handlers := &handlers.BaseRouteHandlers{}
	r.Get("/", handlers.GetLogs)
	r.Get("/ws/live", handlers.LiveWsHandler)
}
