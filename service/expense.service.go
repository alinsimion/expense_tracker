package service

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/model"
)

func NewExpenseService(db db.DB) *ExpenseService {
	return &ExpenseService{

		ExpenseDB: db,
	}
}

type ExpenseService struct {
	ExpenseDB db.DB
}

func (es *ExpenseService) GetExpense(id string) model.Expense {
	return es.ExpenseDB.Expenses[0]
}

func (es *ExpenseService) GetExpenses(skip int, limit int) []model.Expense {
	maxLength := len(es.ExpenseDB.Expenses)

	end := skip + limit

	if end > maxLength {
		end = maxLength
	}

	return es.ExpenseDB.Expenses[skip:end]
}

func (es *ExpenseService) AddExpense(e model.Expense) {
	es.ExpenseDB.Expenses = append(es.ExpenseDB.Expenses, e)
}

func (es *ExpenseService) UpdateExpense(id string, e model.Expense) *model.Expense {
	for i := 0; i < len(es.ExpenseDB.Expenses); i++ {
		if es.ExpenseDB.Expenses[i].Description == id {
			es.ExpenseDB.Expenses[i] = e
			return &e
		}
	}
	return nil
}

// Stats Functions

func (es *ExpenseService) GetCurrentBalance() float64 {
	sum := 0.0

	for i := 0; i < len(es.ExpenseDB.Expenses); i++ {
		sum += es.ExpenseDB.Expenses[i].Amount
	}
	return sum
}

func (es *ExpenseService) GetLargestExpense() model.Expense {
	var e model.Expense

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if tempExpense.Amount > e.Amount {
			e = tempExpense
		}
	}
	return e
}

func (es *ExpenseService) GetAmountByCategory() ([]string, []float64) {
	categories := make(map[string]float64)

	var categoryNames []string
	var categoryValues []float64

	for _, tempExpense := range es.ExpenseDB.Expenses {
		categories[tempExpense.Category] += tempExpense.Amount

		if !slices.Contains(categoryNames, tempExpense.Category) {
			categoryNames = append(categoryNames, tempExpense.Category)
		}
	}

	slices.Sort(categoryNames)

	for _, cName := range categoryNames {
		categoryValues = append(categoryValues, categories[cName])
	}

	return categoryNames, categoryValues
}

func (es *ExpenseService) GetAmountByMonth() ([]string, []float64) {

	months := map[string]float64{
		"January":   0,
		"February":  0,
		"March":     0,
		"April":     0,
		"May":       0,
		"June":      0,
		"July":      0,
		"August":    0,
		"September": 0,
		"October":   0,
		"November":  0,
		"December":  0,
	}

	for _, tempExpense := range es.ExpenseDB.Expenses {
		months[tempExpense.Date.Month().String()] += tempExpense.Amount
	}

	monthIndexes := map[string]float64{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}

	monthNames := []string{"January", "February", "March", "April",
		"May", "June", "July", "August",
		"September", "October", "November", "December"}

	sort.Slice(monthNames, func(i, j int) bool {
		return monthIndexes[monthNames[i]] < monthIndexes[monthNames[j]]
	})

	var monthValues []float64

	for _, mName := range monthNames {
		monthValues = append(monthValues, months[mName])
	}

	return monthNames, monthValues
}

func (es *ExpenseService) GetAmountByDay() map[string]float64 {
	days := make(map[string]float64)

	for _, tempExpense := range es.ExpenseDB.Expenses {
		day := tempExpense.Date.Day()
		month := tempExpense.Date.Month()
		year := tempExpense.Date.Year()

		key := fmt.Sprintf("%02d/%02d/%04d", day, month, year)

		days[key] += tempExpense.Amount
	}
	return days
}

func (es *ExpenseService) GetLongestStreakWithoutExpense() int {

	var dates []time.Time

	for _, tempExpense := range es.ExpenseDB.Expenses {
		dates = append(dates, tempExpense.Date)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	maxStreak := 0
	startDate := dates[0]

	for i := 1; i < len(dates); i++ {

		tempStreak := dates[i].Day() - startDate.Day()

		if tempStreak > maxStreak {
			maxStreak = tempStreak
		}

		startDate = dates[i]
	}

	return maxStreak
}
