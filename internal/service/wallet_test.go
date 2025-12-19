package service

import (
	"context"
	"errors"
	"testing"
	"wallet-test-api/internal/models"

	"github.com/google/uuid"
)

type ManualMockRepository struct {
	getBalanceFn    func(id uuid.UUID) (float64, error)
	createWalletFn  func(id uuid.UUID, balance float64) error
	updateBalanceFn func(id uuid.UUID, balance float64) error
	saveOperationFn func(id uuid.UUID, op string, amount float64) error
}

func (m *ManualMockRepository) WithTransaction(ctx context.Context, operation func(tx any) error) error {
	return operation(nil)
}

func (m *ManualMockRepository) GetBalanceForUpdate(ctx context.Context, tx any, id uuid.UUID) (float64, error) {
	return m.getBalanceFn(id)
}

func (m *ManualMockRepository) CreateWallet(ctx context.Context, id uuid.UUID, balance float64) error {
	return m.createWalletFn(id, balance)
}

func (m *ManualMockRepository) UpdateBalance(ctx context.Context, tx any, id uuid.UUID, balance float64) error {
	return m.updateBalanceFn(id, balance)
}

func (m *ManualMockRepository) SaveOperation(ctx context.Context, tx any, id uuid.UUID, op string, amount float64) error {
	return m.saveOperationFn(id, op, amount)
}

func (m *ManualMockRepository) DeleteWallet(ctx context.Context, id uuid.UUID) error { return nil }

func TestWalletService_ProcessOperation(t *testing.T) {
	walletID := uuid.New()

	tests := []struct {
		name          string
		opType        string
		amount        float64
		setupMock     func(m *ManualMockRepository)
		expectedError error
	}{
		{
			name:   "Успешный депозит (новый кошелек)",
			opType: models.Deposit,
			amount: 1000,
			setupMock: func(m *ManualMockRepository) {
				m.getBalanceFn = func(id uuid.UUID) (float64, error) { return 0, ErrValletNotFound }
				m.createWalletFn = func(id uuid.UUID, b float64) error { return nil }
				m.updateBalanceFn = func(id uuid.UUID, b float64) error { return nil }
				m.saveOperationFn = func(id uuid.UUID, o string, a float64) error { return nil }
			},
			expectedError: nil,
		},
		{
			name:   "Ошибка: недостаточно средств",
			opType: models.Withdrawal,
			amount: 500,
			setupMock: func(m *ManualMockRepository) {
				m.getBalanceFn = func(id uuid.UUID) (float64, error) { return 100, nil }
			},
			expectedError: ErrLowBalance,
		},
		{
			name:   "Ошибка: неверный тип операции",
			opType: "UNKNOWN",
			amount: 10,
			setupMock: func(m *ManualMockRepository) {
				m.getBalanceFn = func(id uuid.UUID) (float64, error) { return 100, nil }
			},
			expectedError: ErrOperationType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ManualMockRepository{}
			tt.setupMock(repo)
			svc := NewWalletService(repo)

			err := svc.ProcessOperation(context.Background(), walletID, tt.opType, tt.amount)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("ожидалась ошибка %v, получена %v", tt.expectedError, err)
			}
		})
	}
}
