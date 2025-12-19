package repository

import (
	"context"
	"fmt"
	"wallet-test-api/internal/service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (repo *PostgresRepo) WithTransaction(ctx context.Context, operation func(tx any) error) error {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	err = operation(tx)
	return err
}

func (repo *PostgresRepo) CreateWallet(ctx context.Context, id uuid.UUID, initialBalance float64) error {
	const query = `INSERT INTO wallets ("valletId", balance) VALUES ($1, $2)`
	_, err := repo.pool.Exec(ctx, query, id, initialBalance)
	if err != nil {
		return fmt.Errorf("failed to create vallet: %w", err)
	}
	return nil
}

func (repo *PostgresRepo) DeleteWallet(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM wallets WHERE "valletId" = $1`
	_, err := repo.pool.Exec(ctx, query, id)
	return err
}

func (repo *PostgresRepo) UpdateBalance(ctx context.Context, transaction any, id uuid.UUID, newBalance float64) error {
	transactionNew, ok := transaction.(pgx.Tx)
	if !ok {
		return fmt.Errorf("internal error: transaction object is invalid")
	}

	const query = `UPDATE wallets SET balance = $1, updated_at = now() WHERE "valletId" = $2`
	_, err := transactionNew.Exec(ctx, query, newBalance, id)
	return err
}

func (repo *PostgresRepo) SaveOperation(ctx context.Context, transaction any, id uuid.UUID, opType string, amount float64) error {
	transactionNew, ok := transaction.(pgx.Tx)
	if !ok {
		return fmt.Errorf("internal error: transaction object is invalid")
	}
	const query = `
		INSERT INTO operations ("valletId", "operationType", "amount") 
		VALUES ($1, $2, $3)`
	_, err := transactionNew.Exec(ctx, query, id, opType, amount)
	return err
}

func (repo *PostgresRepo) GetBalanceForUpdate(ctx context.Context, transaction any, id uuid.UUID) (float64, error) {
	var balance float64
	transactionNew, ok := transaction.(pgx.Tx)
	if ok {
		const query = `SELECT balance FROM wallets WHERE "valletId" = $1 FOR UPDATE`
		err := transactionNew.QueryRow(ctx, query, id).Scan(&balance)
		if err != nil {
			if err == pgx.ErrNoRows {
				return 0, service.ErrValletNotFound
			}
			return 0, err
		}
	} else {
		const query = `SELECT balance FROM wallets WHERE "valletId" = $1`
		err := repo.pool.QueryRow(ctx, query, id).Scan(&balance)
		if err != nil {
			if err == pgx.ErrNoRows {
				return 0, service.ErrValletNotFound
			}
			return 0, err
		}
	}

	return balance, nil
}
