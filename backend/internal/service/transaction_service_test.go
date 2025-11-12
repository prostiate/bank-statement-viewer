package service

import (
	"strings"
	"testing"

	"github.com/prostiate/bank-statement-viewer/internal/model"
	"github.com/prostiate/bank-statement-viewer/internal/repository"
)

func TestProcessCSV(t *testing.T) {
	tests := []struct {
		name        string
		csvContent  string
		wantErr     bool
		errContains string
		wantCount   int
	}{
		{
			name: "valid CSV",
			csvContent: "1624507883,JOHN DOE,DEBIT,250000,SUCCESS,restaurant\n" +
				"1624615065,E-COMMERCE B,DEBIT,150000,PENDING,clothes",
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:        "expected 6 fields, got 7",
			csvContent:  "1624507883,JOHN DOE,DEBIT,250000,SUCCESS,restaurant,08123456",
			wantErr:     true,
			errContains: "expected 6 fields, got 7",
		},
		{
			name:        "expected 6 fields, got 5",
			csvContent:  "1624507883,JOHN DOE,DEBIT,250000,SUCCESS",
			wantErr:     true,
			errContains: "expected 6 fields, got 5",
		},
		{
			name:        "invalid timestamp",
			csvContent:  "2023-06-20T12:34:56,JOHN DOE,DEBIT,250000,SUCCESS,restaurant",
			wantErr:     true,
			errContains: "invalid timestamp",
		},
		{
			name:        "invalid type",
			csvContent:  "1624507883,JOHN DOE,DEBITT,250000,SUCCEED,restaurant",
			wantErr:     true,
			errContains: "invalid type",
		},
		{
			name:        "invalid status",
			csvContent:  "1624507883,JOHN DOE,DEBIT,250000,SUCCEED,restaurant",
			wantErr:     true,
			errContains: "invalid status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.NewTransactionRepository()
			svc := NewTransactionService(repo)

			reader := strings.NewReader(tt.csvContent)
			_, err := svc.ProcessCSV(reader)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ProcessCSV() expected error but got nil")
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ProcessCSV() error = %v, want %v", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ProcessCSV() error = %v, want nil", err)
				return
			}
			transactions := repo.GetAll()
			if len(transactions) != tt.wantCount {
				t.Errorf("ProcessCSV() got %d transactions, want %d", len(transactions), tt.wantCount)
			}
		})
	}
}

func TestGetBalance(t *testing.T) {
	repo := repository.NewTransactionRepository()
	svc := NewTransactionService(repo)

	t.Run("success transactions", func(t *testing.T) {
		repo.Store([]model.Transaction{
			{Type: model.TypeCredit, Amount: 1000, Status: model.StatusSuccess},
			{Type: model.TypeDebit, Amount: 500, Status: model.StatusSuccess},
		})

		balance, err := svc.GetBalance()
		if err != nil {
			t.Errorf("GetBalance() error = %v, want nil", err)
			return
		}
		if balance != 500 {
			t.Errorf("GetBalance() got %f, want 500", balance)
		}
	})

	t.Run("failed transactions", func(t *testing.T) {
		repo.Store([]model.Transaction{
			{Type: model.TypeCredit, Amount: 1000, Status: model.StatusFailed},
			{Type: model.TypeDebit, Amount: 500, Status: model.StatusPending},
		})

		balance, err := svc.GetBalance()
		if err != nil {
			t.Errorf("GetBalance() error = %v, want nil", err)
			return
		}
		if balance != 0 {
			t.Errorf("GetBalance() got %f, want 0", balance)
		}
	})
}
