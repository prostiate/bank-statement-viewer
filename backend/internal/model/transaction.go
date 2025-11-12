package model

import "time"

type TransactionType string
type TransactionStatus string

const (
	TypeDebit  TransactionType = "DEBIT"
	TypeCredit TransactionType = "CREDIT"
)

const (
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusFailed  TransactionStatus = "FAILED"
	StatusPending TransactionStatus = "PENDING"
)

type Transaction struct {
	Timestamp   time.Time         `json:"timestamp"`
	Name        string            `json:"name"`
	Type        TransactionType   `json:"type"`
	Amount      float64           `json:"amount"`
	Status      TransactionStatus `json:"status"`
	Description string            `json:"description"`
}
