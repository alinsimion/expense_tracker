package model

import "time"

var INCOME = "Income"
var EXPENSE = "Expense"

type Expense struct {
	Id          string    `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	Currency    string    `json:"currency"`
	Type        string    `json:"type"` // expense or income
}
