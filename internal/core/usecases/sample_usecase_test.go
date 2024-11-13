package usecases

import (
	"context"
	"testing"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/mocks/mockpiyographql"
	"github.com/golang/mock/gomock"
)

func TestSampleUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mockpiyographql.NewMockClient(ctrl)
	target := NewSampleUsecase(logger.NewLogger(&config.AppConfig{}), mockClient) // fixme test cfg

	t.Run("get sample", func(t *testing.T) {
		ID := "123"
		// モックの振る舞いを設定
		mockClient.EXPECT().GetSample(context.Background(), ID).
			Return(&models.Sample{ID: "123", StringVal: "Test Sample"}, nil)

		// テストケースを実行
		sample, err := target.Get(context.Background(), ID)

		assert.NoError(t, err)
		assert.Equal(t, "123", sample.ID)
	})

	t.Run("get sample 2", func(t *testing.T) {
		ID := "aaa"
		// モックの振る舞いを設定
		mockClient.EXPECT().GetSample(context.Background(), ID).
			Return(&models.Sample{ID: "aaa", StringVal: "Test Sample"}, nil)

		// テストケースを実行
		sample, err := target.Get(context.Background(), ID)

		assert.NoError(t, err)
		assert.Equal(t, "aaa", sample.ID)
	})
}
