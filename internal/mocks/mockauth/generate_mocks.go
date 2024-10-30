package mockauth

//go:generate mockgen -package=mockauth -destination=./mock_auth.go github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/auth TokenService
