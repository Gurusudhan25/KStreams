package transaction

import (
	"fmt"
	"time"
)

type TransactionService struct {
	transactions []Transaction
}

func NewTransactionService() *TransactionService {
	return &TransactionService{
		transactions: make([]Transaction, 0),
	}
}

func (ts *TransactionService) AddTransaction(t Transaction) {
	ts.transactions = append(ts.transactions, t)
}

func ParseTransaction(record []string) (Transaction, error) {
	var transaction Transaction
	layout := "Mon Jan 02 15:04:05 MST 2006"

	transaction.UserId = record[0]
	transaction.TransactionId = record[1]
	transactionTime, err := time.Parse(layout, record[2])
	if err != nil {
		return transaction, err
	}
	transaction.TransactionTime = transactionTime
	transaction.ItemCode = record[3]
	transaction.ItemDescription = record[4]
	transaction.NumberOfItemsPurchased = atoi(record[5])
	transaction.CostPerItem = atof(record[6])
	transaction.Country = record[7]

	return transaction, nil
}

func atoi(s string) int {
	i := 0
	fmt.Sscanf(s, "%d", &i)
	return i
}

func atof(s string) float64 {
	f := 0.0
	fmt.Sscanf(s, "%f", &f)
	return f
}
