package handlers

import (
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type HealthcheckHandler struct {
	logger     logger.Logger
	JSONWriter *presenter.JSONWriter
}

func NewHealthcheckHandler(logger logger.Logger, JSONWriter *presenter.JSONWriter) *HealthcheckHandler {
	return &HealthcheckHandler{
		logger:     logger,
		JSONWriter: JSONWriter,
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
	//ctx := r.Context()
	//tracer := otel.Tracer("health-check-handler")
	//_, span := tracer.Start(ctx, "health-check")
	//defer span.End()
	//
	//// スパンにカスタム属性を追加
	//span.SetAttributes(attribute.String("custom.attribute", "test-value"))

	span, ctx := tracer.StartSpanFromContext(r.Context(), "test.operation",
		tracer.SpanType("web"),
		tracer.ResourceName("/test"),
	)
	defer span.Finish()
	span.SetTag("custom.tag", "test-value")

	h.logger.InfoContext(ctx, "Healthcheck ok")

	res := map[string]string{"message": "healthcheck ok"}
	h.JSONWriter.Write(ctx, w, res)
}
