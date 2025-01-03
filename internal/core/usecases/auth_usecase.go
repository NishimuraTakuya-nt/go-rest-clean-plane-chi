package usecases

import (
	"context"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/services"
)

type AuthUsecase interface {
	Login(ctx context.Context, userID string, roles []string) (string, error)
	Authenticate(ctx context.Context, tokenString string) (*models.User, error)
}

type authUsecase struct {
	tokenService services.TokenService
}

func NewAuthUsecase(tokenService services.TokenService) AuthUsecase {
	return &authUsecase{
		tokenService: tokenService,
	}
}

func (uc *authUsecase) Login(ctx context.Context, userID string, roles []string) (string, error) {
	return uc.tokenService.GenerateToken(ctx, userID, roles)
}

func (uc *authUsecase) Authenticate(ctx context.Context, tokenString string) (*models.User, error) {
	claims, err := uc.tokenService.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	return &models.User{
		ID:    claims.UserID,
		Roles: claims.Roles,
	}, nil
}
