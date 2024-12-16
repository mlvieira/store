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

// OrderRepository defines methods to interact with order data.
type OrderRepository interface {
	InsertOrder(ctx context.Context, order models.Order) (int, error)
}

// CustomerRepository defines methods to interact with customer data.
type CustomerRepository interface {
	InsertCustomer(ctx context.Context, customer models.Customer) (int, error)
}

// Repositories aggregates repository interfaces.
type Repositories struct {
	Widget      WidgetRepository
	Transaction TransactionRepository
	Order       OrderRepository
	Customer    CustomerRepository
}

// NewRepositories initializes repositories with a database connection.
func NewRepositories(conn *sql.DB) *Repositories {
	return &Repositories{
		Widget:      NewWidgetRepository(conn),
		Transaction: NewTransactionRepository(conn),
		Order:       NewOrderRepository(conn),
		Customer:    NewCustomerRepository(conn),
	}
}
