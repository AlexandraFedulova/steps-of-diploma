package handler

import (
	"net/http"

	"github.com/mephi/sociometry/pkg/httputil"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "sociometry-backend",
	})
}
