package http

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/thoas/stats"
)

type adminHandler struct {
	Stats *stats.Stats
}

func newAdminHandler(s *stats.Stats) http.Handler {
	return handlers.MethodHandler{
		"GET": &adminHandler{
			Stats: s,
		},
	}
}

func (h *adminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s := h.Stats.Data()
	b, _ := json.Marshal(s)

	_, err := w.Write(b)
	if err != nil {
		log.Error(err)
	}
}
