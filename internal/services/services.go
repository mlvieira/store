package services

import (
	"github.com/mlvieira/store/internal/repository"
)

// Services contains all application service instances.
type Services struct {
	CustomerService    *CustomerService
	OrderService       *OrderService
	TransactionService *TransactionService
}

// NewServices initializes and returns all application services.
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		CustomerService:    NewCustomerService(repos.Customer),
		OrderService:       NewOrderService(repos.Order, repos.Customer, repos.Transaction),
		TransactionService: NewTransactionService(repos.Transaction),
	}
}
