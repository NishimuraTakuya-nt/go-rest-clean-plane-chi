package mockrepository

//go:generate mockgen -package=mockrepository -destination=./mock_reposiotry.go github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/secondary/repository OrderRepository,ProductRepository
