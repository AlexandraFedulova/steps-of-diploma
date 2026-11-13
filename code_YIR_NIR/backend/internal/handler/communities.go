package handler

import (
	"bytes"
	"net/http"

	"github.com/mephi/sociometry/pkg/graph"
	"github.com/mephi/sociometry/pkg/httputil"
)

type CommunitiesHandler struct{}

func NewCommunitiesHandler() *CommunitiesHandler {
	return &CommunitiesHandler{}
}

func (h *CommunitiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	data, header, err := readUploadFile(w, r)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid multipart form or missing file field")
		return
	}

	g, err := graph.ParseCSV(bytes.NewReader(data))
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	cliques := g.CliqueCommunities()

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"filename": header.Filename,
		"size":     len(data),
		"method":   "clique",
		"cliques":  cliques,
	})
}
