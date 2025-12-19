package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

type MockWalletService struct {
	processOpFn func(id uuid.UUID, opType string, amount float64) error
}

func (m *MockWalletService) ProcessOperation(ctx context.Context, id uuid.UUID, op string, amt float64) error {
	return m.processOpFn(id, op, amt)
}

func (m *MockWalletService) GetBalance(ctx context.Context, id uuid.UUID) (float64, error) {
	return 0, nil
}

func TestWalletHandler_HandleOperation(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "Успешная операция",
			requestBody:    `{"valletId":"550e8400-e29b-41d4-a716-446655440000","operationType":"DEPOSIT","amount":100}`,
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Ошибка: Кривой JSON",
			requestBody:    `{"valletId":"invalid-uuid",,,}`,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Ошибка: Сервис вернул ошибку",
			requestBody:    `{"valletId":"550e8400-e29b-41d4-a716-446655440000","operationType":"WITHDRAW","amount":1000}`,
			mockReturnErr:  errors.New("insufficient funds"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &MockWalletService{
				processOpFn: func(id uuid.UUID, op string, amt float64) error {
					return tt.mockReturnErr
				},
			}

			h := NewWalletHandler(mockSvc)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(tt.requestBody))
			rec := httptest.NewRecorder()

			h.HandleOperation(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получили %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
