package services

import (
	"context"

	"github.com/mlvieira/store/internal/models"
	"github.com/mlvieira/store/internal/repository"
)

type CustomerService struct {
	repo repository.CustomerRepository
}

// NewCustomerService initializes a new CustomerService instance.
func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

// SaveCustomer saves a customer and returns the ID.
func (s *CustomerService) SaveCustomer(ctx context.Context, customer models.Customer) (int, error) {
	return s.repo.InsertCustomer(ctx, customer)
}
