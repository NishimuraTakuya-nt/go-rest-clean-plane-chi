package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

type HealthcheckHandler struct {
	log logger.Logger
}

func NewHealthcheckHandler(log logger.Logger) *HealthcheckHandler {
	return &HealthcheckHandler{
		log: log,
	}
}

// Get godoc
// @Summary Health check endpoint
// @Description Get the health status of the API
// @Tags healthcheck
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} response.ErrorResponse
// @Router /healthcheck [get]
func (h *HealthcheckHandler) Get(w http.ResponseWriter, r *http.Request) {
	// healthcheck
	h.log.InfoContext(r.Context(), "Healthcheck ok")
	// nolint:errcheck
	json.NewEncoder(w).Encode(map[string]string{"message": "healthcheck ok"})
}
