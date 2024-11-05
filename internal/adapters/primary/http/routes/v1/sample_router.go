package v1

import (
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/handlers"
	"github.com/go-chi/chi/v5"
)

type SampleRouter struct {
	Handler http.Handler
}

func NewSampleRouter(sampleHandler *handlers.SampleHandler) *SampleRouter {
	r := chi.NewRouter()

	r.Get("/", sampleHandler.List)
	r.Post("/", sampleHandler.Create)

	// ID指定の操作をグループ化
	r.Route("/{sampleID}", func(r chi.Router) {
		r.Get("/", sampleHandler.Get)
		r.Put("/", sampleHandler.Update)
		r.Delete("/", sampleHandler.Delete)

		// ネストされたリソース
		r.Route("/profile", func(r chi.Router) {
			r.Get("/", sampleHandler.GetSampleProfile)
			r.Put("/", sampleHandler.UpdateSampleProfile)
		})

	})

	return &SampleRouter{Handler: r}
}
