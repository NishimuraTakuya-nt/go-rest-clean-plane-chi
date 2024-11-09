package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx := r.Context()
	tracer := otel.Tracer("health-check-handler")
	_, span := tracer.Start(ctx, "health-check")
	defer span.End()

	// スパンにカスタム属性を追加
	span.SetAttributes(attribute.String("custom.attribute", "test-value"))

	// healthcheck
	h.log.InfoContext(ctx, "Healthcheck ok")
	// nolint:errcheck
	json.NewEncoder(w).Encode(map[string]string{"message": "healthcheck ok"})
}
