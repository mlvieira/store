package services

import (
	"context"

	"github.com/mlvieira/store/internal/models"
	"github.com/mlvieira/store/internal/repository"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

// NewTransactionService initializes a new TransactionService instance.
func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// SaveTransaction saves a transaction and returns the ID.
func (s *TransactionService) SaveTransaction(ctx context.Context, txn models.Transaction) (int, error) {
	return s.repo.InsertTransaction(ctx, txn)
}
