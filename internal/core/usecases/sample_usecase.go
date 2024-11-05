package usecases

import (
	"context"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/secondary/piyographql"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

type SampleUsecase interface {
	Get(ctx context.Context, ID string) (*models.Sample, error)
	List(ctx context.Context, offset, limit *int) ([]models.Sample, error)
}

type sampleUsecase struct {
	log           logger.Logger
	graphqlClient piyographql.Client
}

func NewSampleUsecase(log logger.Logger, client piyographql.Client) SampleUsecase {
	return &sampleUsecase{
		log:           log,
		graphqlClient: client,
	}
}

func (uc *sampleUsecase) Get(ctx context.Context, ID string) (*models.Sample, error) {
	// todo trace log
	return uc.graphqlClient.GetSample(ctx, ID)
}

func (uc *sampleUsecase) List(ctx context.Context, offset, limit *int) ([]models.Sample, error) {
	return uc.graphqlClient.ListSample(ctx, offset, limit)
}
