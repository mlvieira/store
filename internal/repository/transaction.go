package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mlvieira/store/internal/models"
)

// transactionRepo handles database operations for transactions.
type transactionRepo struct {
	db *sql.DB
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

// InsertTransaction inserts a new transaction into the database.
func (r *transactionRepo) InsertTransaction(ctx context.Context, txn models.Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt := `
		INSERT INTO transactions 
		(amount, currency, last_four, bank_return_code, 
		 transaction_status_id, created_at, updated_at, 
		 expiry_month, expiry_year, payment_intent, payment_method)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return int(id), nil
}
