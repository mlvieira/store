package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/mlvieira/store/internal/models"
)

// customerRepo handles database operations for customer.
type customerRepo struct {
	db *sql.DB
}

// NewCustomerRepository creates a new customerRepository
func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepo{db: db}
}

// InsertCustomer inserts a new customer into the database.
func (r *customerRepo) InsertCustomer(ctx context.Context, customer models.Customer) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	stmt := `
		INSERT INTO customers 
		(first_name, last_name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(ctx, stmt,
		customer.FirstName,
		customer.LastName,
		customer.Email,
		time.Now(),
		time.Now(),
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
