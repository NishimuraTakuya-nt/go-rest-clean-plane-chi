package response

import (
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
)

// SampleResponse はサンプルのレスポンスを表す構造体です
// @Description Sample information
type SampleResponse struct {
	ID        string    `json:"id"`
	StringVal string    `json:"string_val"`
	IntVal    int       `json:"int_val"`
	ArrayVal  []string  `json:"array_val"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToSampleResponse はドメインモデルからレスポンスモデルへの変換を行います
func ToSampleResponse(s *models.Sample) SampleResponse {
	return SampleResponse{
		ID:        s.ID,
		StringVal: s.StringVal,
		IntVal:    s.IntVal,
		ArrayVal:  s.ArrayVal,
		Email:     s.Email,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ListSampleResponse は複数サンプルを返すためのレスポンス構造体です
// @Description Sample list information
type ListSampleResponse struct {
	Samples    []SampleResponse `json:"samples"`
	TotalCount int              `json:"total_count"`
	Offset     *int             `json:"offset"`
	Limit      *int             `json:"limit"`
}

// ToListSampleResponse は複数のサンプルモデルを変換します
func ToListSampleResponse(models []models.Sample, offset, limit *int) *ListSampleResponse {
	samples := make([]SampleResponse, len(models))
	for i, model := range models {
		samples[i] = ToSampleResponse(&model)
	}

	return &ListSampleResponse{
		Samples:    samples,
		TotalCount: len(samples),
		Offset:     offset,
		Limit:      limit,
	}
}
