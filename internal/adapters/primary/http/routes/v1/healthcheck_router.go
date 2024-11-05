package v1

import (
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/handlers"
	"github.com/go-chi/chi/v5"
)

type HealthcheckRouter struct {
	Handler http.Handler
}

func NewHealthcheckRouter(healthcheckHandler *handlers.HealthcheckHandler) *HealthcheckRouter {
	r := chi.NewRouter()
	r.Get("/", healthcheckHandler.Get)

	return &HealthcheckRouter{Handler: r}
}
