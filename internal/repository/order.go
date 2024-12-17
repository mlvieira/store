package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mlvieira/store/internal/models"
)

// orderRepo handles database operations for order.
type orderRepo struct {
	db *sql.DB
}

// NewOrderRepository creates a new orderRepository
func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepo{db: db}
}

// InsertOrder inserts a new order into the database.
func (r *orderRepo) InsertOrder(ctx context.Context, order models.Order) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt := `
		INSERT INTO orders  
		(widget_id, transaction_id, status_id, quantity, 
		 amount, created_at, updated_at, customer_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
		order.CustomerID,
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
