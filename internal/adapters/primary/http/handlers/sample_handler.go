package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/dto/request"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/dto/response"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/handlers/queryparameter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/apperrors"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/usecases"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type SampleHandler struct {
	log           logger.Logger
	sampleUsecase usecases.SampleUsecase
}

func NewSampleHandler(log logger.Logger, sampleUsecase usecases.SampleUsecase) *SampleHandler {
	return &SampleHandler{
		log:           log,
		sampleUsecase: sampleUsecase,
	}
}

// List godoc
// @Summary List samples
// @Description Get a list of samples with pagination
// @Tags samples
// @Accept  json
// @Produce  json
// @Param offset query int false "Offset for pagination" default(0) minimum(0)
// @Param limit query int false "Limit for pagination" default(100) minimum(1) maximum(100)
// @Security ApiKeyAuth
// @Success 200 {object} response.ListSampleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /samples [get]
func (h *SampleHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	p := queryparameter.NewOffsetLimitParams(r)
	if err := validator.Validate(p); err != nil {
		h.log.ErrorContext(ctx, "Invalid parameters", "error", err)
		writeError(w, err)
		return
	}

	// サンプルリストの取得
	samples, err := h.sampleUsecase.List(ctx, p.Offset, p.Limit)
	if err != nil {
		h.log.ErrorContext(ctx, "Failed to get sample list", "error", err)
		writeError(w, err)
		return
	}

	res := response.ToListSampleResponse(samples, p.Offset, p.Limit)

	writeJSONResponse(ctx, w, res)
}

// Create godoc
// @Summary Sample create
// @Description Create a new sample
// @Tags samples
// @Accept json
// @Produce json
// @Param request body request.SampleRequest true "Sample information"
// @Security ApiKeyAuth
// @Success 200 {object} response.SampleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /samples [post]
func (h *SampleHandler) Create(w http.ResponseWriter, r *http.Request) {
	// サンプル作成処理
	ctx := r.Context()

	var req request.SampleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.ErrorContext(ctx, "Failed to decode sample request", "error", err)
		writeError(w, apperrors.NewBadRequestError("Invalid request body", err))
		return
	}

	if validationErrors := validator.Validate(req); validationErrors != nil {
		writeError(w, validationErrors)
		return
	}

	res := response.SampleResponse{
		ID:        "123",
		StringVal: req.StringVal,
		IntVal:    req.IntVal,
		ArrayVal:  req.ArrayVal,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	writeJSONResponse(ctx, w, res)
}

// Get godoc
// @Summary Get a sample by ID
// @Description Get details of a sample
// @Tags samples
// @Accept  json
// @Produce  json
// @Param id path string true "Sample ID"
// @Security ApiKeyAuth
// @Success 200 {object} response.SampleResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /samples/{id} [get]
func (h *SampleHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ID := chi.URLParam(r, "id")
	if ID == "" {
		h.log.ErrorContext(ctx, "ID is required")
		writeError(w, apperrors.NewBadRequestError("ID is required", nil))
		return
	}
	if err := validator.ValidateVar(ID, "sampleId", "path parameter"); err != nil {
		h.log.WarnContext(ctx, "Invalid sample ID format", "id", ID)
		writeError(w, err)
		return
	}

	sample, err := h.sampleUsecase.Get(ctx, ID)
	if err != nil {
		h.log.ErrorContext(ctx, "Failed to get sample", "error", err)
		writeError(w, err)
		return
	}

	res := response.ToSampleResponse(sample)

	writeJSONResponse(ctx, w, res)
}

func (h *SampleHandler) Update(w http.ResponseWriter, _ *http.Request) {
	// 更新処理
	// nolint:errcheck
	json.NewEncoder(w).Encode(map[string]string{"message": "Update sample"})
}

func (h *SampleHandler) Delete(w http.ResponseWriter, _ *http.Request) {
	// 削除処理
	// nolint:errcheck
	json.NewEncoder(w).Encode(map[string]string{"message": "Delete sample"})
}

func (h *SampleHandler) GetSampleProfile(_ http.ResponseWriter, _ *http.Request) {
}

func (h *SampleHandler) UpdateSampleProfile(_ http.ResponseWriter, _ *http.Request) {
}
