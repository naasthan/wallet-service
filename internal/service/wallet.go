package service

import (
	"context"
	"errors"
	"fmt"
	"wallet-test-api/internal/models"

	"github.com/google/uuid"
)

var (
	ErrValletNotFound = errors.New("vallet not found")
	ErrLowBalance     = errors.New("low balance")
	ErrOperationType  = errors.New("operation type")
)

type WalletRepository interface {
	WithTransaction(ctx context.Context, operation func(tx any) error) error
	CreateWallet(ctx context.Context, id uuid.UUID, initialBalance float64) error
	GetBalanceForUpdate(ctx context.Context, transaction any, id uuid.UUID) (float64, error)
	DeleteWallet(ctx context.Context, id uuid.UUID) error
	UpdateBalance(ctx context.Context, transaction any, id uuid.UUID, newBalance float64) error
	SaveOperation(ctx context.Context, transaction any, id uuid.UUID, opType string, amount float64) error
}

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (service *WalletService) ProcessOperation(ctx context.Context, walletID uuid.UUID, opType string, amount float64) error {
	return service.repo.WithTransaction(ctx, func(tx any) error {
		currentBalance, err := service.repo.GetBalanceForUpdate(ctx, tx, walletID)

		if err != nil {
			fmt.Printf("Wallet %s not found, checking if deposit...\n", walletID)

			if opType == models.Deposit {
				if createErr := service.repo.CreateWallet(ctx, walletID, 0); createErr != nil {
					fmt.Printf("FAILED to create wallet: %v\n", createErr)
					return createErr
				}
				currentBalance = 0
			} else {
				return err
			}
		}

		var newBalance float64

		switch opType {
		case models.Deposit:
			newBalance = currentBalance + amount
		case models.Withdrawal:
			if currentBalance < amount {
				return ErrLowBalance
			}
			newBalance = currentBalance - amount
		default:
			return ErrOperationType
		}

		if err := service.repo.UpdateBalance(ctx, tx, walletID, newBalance); err != nil {
			return err
		}

		return service.repo.SaveOperation(ctx, tx, walletID, opType, amount)
	})
}

func (s *WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	return s.repo.GetBalanceForUpdate(ctx, nil, walletID)
}
