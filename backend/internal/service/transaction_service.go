package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/prostiate/bank-statement-viewer/internal/model"
	"github.com/prostiate/bank-statement-viewer/internal/repository"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{
		repo: repo,
	}
}

func isEmptyLine(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

func (s *TransactionService) ProcessCSV(file io.Reader) ([]model.Transaction, error) {
	reader := csv.NewReader(file)
	transactions := make([]model.Transaction, 0)

	const maxRecords = 100_000
	lineNum := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		lineNum++

		if err != nil {
			return nil, fmt.Errorf("line %d: failed to read CSV: %w", lineNum, err)
		}

		if lineNum > maxRecords {
			return nil, fmt.Errorf("CSV contains too many records, max allowed: %d", maxRecords)
		}

		if isEmptyLine(record) {
			continue
		}

		if len(record) != 6 {
			return nil, fmt.Errorf("line %d: expected 6 fields, got %d", lineNum, len(record))
		}

		timestampStr := strings.TrimSpace(record[0])
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid timestamp: %w", lineNum, err)
		}
		parsedTimestamp := time.Unix(timestamp, 0)

		typeStr := strings.ToUpper(strings.TrimSpace(record[2]))
		txType := model.TransactionType(typeStr)
		if txType != model.TypeCredit && txType != model.TypeDebit {
			return nil, fmt.Errorf("line %d: invalid type: '%s', must be 'CREDIT' or 'DEBIT'", lineNum, typeStr)
		}

		amountStr := strings.TrimSpace(record[3])
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid amount: %w", lineNum, err)
		}

		statusStr := strings.ToUpper(strings.TrimSpace(record[4]))
		txStatus := model.TransactionStatus(statusStr)
		if txStatus != model.StatusSuccess && txStatus != model.StatusFailed && txStatus != model.StatusPending {
			return nil, fmt.Errorf("line %d: invalid status: '%s', must be 'SUCCESS', 'FAILED' or 'PENDING'", lineNum, statusStr)
		}

		transaction := model.Transaction{
			Timestamp:   parsedTimestamp,
			Name:        strings.TrimSpace(record[1]),
			Type:        txType,
			Amount:      amount,
			Status:      txStatus,
			Description: strings.TrimSpace(record[5]),
		}
		transactions = append(transactions, transaction)
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no valid transactions found in CSV")
	}

	if err := s.repo.Store(transactions); err != nil {
		return nil, fmt.Errorf("failed to store transactions: %w", err)
	}

	return transactions, nil
}

func (s *TransactionService) GetBalance() (float64, error) {
	transactions := s.repo.GetAll()
	var balance float64
	for _, t := range transactions {
		if t.Status != model.StatusSuccess {
			continue
		}
		switch t.Type {
		case model.TypeCredit:
			balance += t.Amount
		case model.TypeDebit:
			balance -= t.Amount
		}
	}
	return balance, nil
}

func (s *TransactionService) GetIssues() ([]model.Transaction, error) {
	transactions := s.repo.GetAll()
	var issues []model.Transaction
	for _, t := range transactions {
		if t.Status == model.StatusFailed || t.Status == model.StatusPending {
			issues = append(issues, t)
		}
	}
	return issues, nil
}
