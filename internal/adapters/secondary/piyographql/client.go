package piyographql

import (
	"context"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

type Client interface {
	GetSample(ctx context.Context, id string) (*models.Sample, error)
	ListSample(ctx context.Context, offset, limit *int) ([]models.Sample, error)
}

type client struct {
	log logger.Logger
	// クライアントの設定など
}

func NewClient(log logger.Logger) Client {
	return &client{
		log: log,
	}
}

func (c *client) GetSample(_ context.Context, ID string) (*models.Sample, error) {
	// ここでは簡易的に固定値を返していますが、
	// 実際には取得する処理を実装します

	return &models.Sample{
		ID:        ID,
		StringVal: "example1",
		IntVal:    123,
		ArrayVal:  []string{"aaa", "bbb", "ccc"},
		Email:     "user@example.com",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}, nil
}

func (c *client) ListSample(ctx context.Context, offset, limit *int) ([]models.Sample, error) {
	c.log.InfoContext(ctx, "client ListSample", "offset", offset, "limit", limit)
	return []models.Sample{
		{
			ID:        "1",
			StringVal: "example1",
			IntVal:    123,
			ArrayVal:  []string{"aaa", "bbb", "ccc"},
		},
		{
			ID:        "2",
			StringVal: "example2",
			IntVal:    124,
			ArrayVal:  []string{"ddd", "eee", "fff"},
		},
	}, nil
}
