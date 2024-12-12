package repository

import (
	"context"

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
