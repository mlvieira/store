package repository

import (
	"context"
	"database/sql"

	"github.com/mlvieira/store/internal/models"
)

type WidgetRepository interface {
	GetWidgetByID(ctx context.Context, id int) (models.Widget, error)
}

type TransactionRepository interface {
	InsertTransaction(ctx context.Context, txn models.Transaction) (int, error)
}

type Repositories struct {
	Widget      WidgetRepository
	Transaction TransactionRepository
}

func NewRepositories(conn *sql.DB) *Repositories {
	return &Repositories{
		Widget:      NewWidgetRepository(conn),
		Transaction: NewTransactionRepository(conn),
	}
}
