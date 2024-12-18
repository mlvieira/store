package services

import (
	"context"

	"github.com/mlvieira/store/internal/models"
	"github.com/mlvieira/store/internal/repository"
)

type OrderService struct {
	order repository.OrderRepository
}

// NewOrderService initializes a new OrderService instance.
func NewOrderService(order repository.OrderRepository) *OrderService {
	return &OrderService{order: order}
}

// PlaceOrder processes an order and returns the order ID.
func (s *OrderService) PlaceOrder(ctx context.Context, order models.Order) (int, error) {
	return s.order.InsertOrder(ctx, order)
}
