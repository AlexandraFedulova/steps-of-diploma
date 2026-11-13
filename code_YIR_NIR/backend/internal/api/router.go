package api

import (
	"net/http"

	"github.com/mephi/sociometry/internal/handler"
	"github.com/mephi/sociometry/internal/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", handler.NewHealthHandler())
	mux.Handle("POST /upload", handler.NewUploadHandler())
	mux.Handle("POST /communities", handler.NewCommunitiesHandler())

	return middleware.CORS(mux)
}
