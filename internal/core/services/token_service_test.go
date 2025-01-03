package services

import (
	"context"
	"testing"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/mocks/mockservice"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenService := mockservice.NewMockTokenService(ctrl)

	testCases := []struct {
		name     string
		userID   string
		roles    []string
		mockFunc func()
		want     string
		wantErr  bool
	}{
		{
			name:   "正常系: トークンが正しく生成される",
			userID: "user123",
			roles:  []string{"user", "admin"},
			mockFunc: func() {
				mockTokenService.EXPECT().
					GenerateToken(context.Background(), "user123", []string{"user", "admin"}).
					Return("mocked.jwt.token", nil)
			},
			want:    "mocked.jwt.token",
			wantErr: false,
		},
		{
			name:   "異常系: エラーが返される",
			userID: "user456",
			roles:  []string{"user"},
			mockFunc: func() {
				mockTokenService.EXPECT().
					GenerateToken(context.Background(), "user456", []string{"user"}).
					Return("", jwt.ErrSignatureInvalid)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			got, err := mockTokenService.GenerateToken(context.Background(), tc.userID, tc.roles)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
