package handlers

import "net/http"

type BaseRouteHandlers struct{}

func (h *BaseRouteHandlers) LiveWsHandler(w http.ResponseWriter, r *http.Request) {}

func (h *BaseRouteHandlers) GetLogs(w http.ResponseWriter, r *http.Request) {}
