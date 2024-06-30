package transaction

import "time"

type Transaction struct {
	UserId                 string
	TransactionId          string
	TransactionTime        time.Time
	ItemCode               string
	ItemDescription        string
	NumberOfItemsPurchased int
	CostPerItem            float64
	Country                string
}
