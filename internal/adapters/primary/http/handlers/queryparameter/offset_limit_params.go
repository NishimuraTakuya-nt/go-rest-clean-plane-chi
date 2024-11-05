package queryparameter

import (
	"net/http"
	"strconv"
)

type OffsetLimitParams struct {
	Offset *int `validate:"omitempty,gte=0"`
	Limit  *int `validate:"omitempty,gte=1,lte=100"`
}

// NewOffsetLimitParams クエリパラメータから OffsetLimitParams を生成します
func NewOffsetLimitParams(r *http.Request) OffsetLimitParams {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")
	params := OffsetLimitParams{}

	if offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			params.Offset = &val
		}
	}
	if limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			params.Limit = &val
		}
	}
	return params
}
