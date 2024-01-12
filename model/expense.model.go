package model

import "time"

type Expense struct {
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Currency    string    `json:"currency"`
	Type        string    `json:"type"` // expense or income
}
