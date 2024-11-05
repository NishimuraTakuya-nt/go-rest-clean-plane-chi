package repository

import (
	"context"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
)

type SampleRepository interface {
	Get(ctx context.Context, id string) (*models.Sample, error)
	List(ctx context.Context, offset, limit *int) ([]*models.Sample, error)
}
