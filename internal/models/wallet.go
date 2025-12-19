package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	Deposit    = "DEPOSIT"
	Withdrawal = "WITHDRAW"
)

type Wallet struct {
	ValletID  uuid.UUID `json:"valletId"`
	Balance   float64   `json:"balance"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WalletRequest struct {
	ValletID      uuid.UUID `json:"valletId"`
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"`
}
