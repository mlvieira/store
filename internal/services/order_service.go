package services

import (
	"context"

	"github.com/mlvieira/store/internal/models"
	"github.com/mlvieira/store/internal/repository"
)

type OrderService struct {
	orderRepo       repository.OrderRepository
	customerRepo    repository.CustomerRepository
	transactionRepo repository.TransactionRepository
}

// NewOrderService initializes a new OrderService instance.
func NewOrderService(
	orderRepo repository.OrderRepository,
	customerRepo repository.CustomerRepository,
	transactionRepo repository.TransactionRepository,
) *OrderService {
	return &OrderService{
		orderRepo:       orderRepo,
		customerRepo:    customerRepo,
		transactionRepo: transactionRepo,
	}
}

// PlaceOrder processes an order and returns the order ID.
func (s *OrderService) PlaceOrder(ctx context.Context, order models.Order) (int, error) {
	return s.orderRepo.InsertOrder(ctx, order)
}
