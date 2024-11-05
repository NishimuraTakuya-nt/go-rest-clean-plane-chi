package v1

import (
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/handlers"
	"github.com/go-chi/chi/v5"
)

type AuthRouter struct {
	Handler http.Handler
}

func NewAuthRouter(authHandler *handlers.AuthHandler) *AuthRouter {
	r := chi.NewRouter()
	r.Post("/login", authHandler.Login)

	return &AuthRouter{Handler: r}
}
