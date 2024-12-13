package repository

import (
	"context"
	"database/sql"

	"github.com/mlvieira/store/internal/models"
)

// WidgetRepository defines methods to interact with widget data.
type WidgetRepository interface {
	GetWidgetByID(ctx context.Context, id int) (models.Widget, error)
}

// TransactionRepository defines methods to interact with transaction data.
type TransactionRepository interface {
	InsertTransaction(ctx context.Context, txn models.Transaction) (int, error)
}

// Repositories aggregates repository interfaces.
type Repositories struct {
	Widget      WidgetRepository
	Transaction TransactionRepository
}

// NewRepositories initializes repositories with a database connection.
func NewRepositories(conn *sql.DB) *Repositories {
	return &Repositories{
		Widget:      NewWidgetRepository(conn),
		Transaction: NewTransactionRepository(conn),
	}
}
