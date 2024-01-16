package model

type SortedMap struct {
	Keys   []string
	Values []float64
}

type FilterFunc func(Expense) bool
