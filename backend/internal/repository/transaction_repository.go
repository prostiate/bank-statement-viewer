package repository

import (
	"sync"

	"github.com/prostiate/bank-statement-viewer/internal/model"
)

type TransactionRepository interface {
	Store(transactions []model.Transaction) error
	GetAll() []model.Transaction
}

type MuTransactionRepository struct {
	transactions []model.Transaction
	mu           sync.Mutex
}

func NewTransactionRepository() TransactionRepository {
	return &MuTransactionRepository{
		transactions: make([]model.Transaction, 0),
	}
}

func (r *MuTransactionRepository) Store(transactions []model.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.transactions = transactions
	return nil
}

func (r *MuTransactionRepository) GetAll() []model.Transaction {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.transactions
}
