package model

import "time"

const (
	INCOME  = "Income"
	EXPENSE = "Expense"
	BOTH    = "Both"
)

type Expense struct {
	Id          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Currency    string    `json:"currency"`
	Type        string    `json:"type"` // expense or income
	Uid         int64     `json:"uid"`
}
