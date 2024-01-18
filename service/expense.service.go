package service

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/model"
)

type ExpenseServiceInterface interface {
	GetExpense(id string) *model.Expense
	GetExpenses(skip int, limit int) []model.Expense
	AddExpense(e model.Expense)
	UpdateExpense(id string, e model.Expense) *model.Expense
	DeleteExpense(id string) error
	GetCurrentBalance(filter model.FilterFunc) float64
	GetLargestExpense(filter model.FilterFunc) model.Expense
	GetExpensesByCategory(filter model.FilterFunc) ([]string, []float64)
	GetExpensesByMonth(filter model.FilterFunc) ([]string, []float64)
	GetExpensesByDay(filter model.FilterFunc) map[string]float64
	GetLongestStreakWithoutExpense(filter model.FilterFunc) int
	GetCountsByCategory(filter model.FilterFunc) ([]string, []float64)
	GetCurrentIncomes(filter model.FilterFunc) float64
	GetCurrentExpenses(filter model.FilterFunc) float64
}

func NewExpenseService(db db.DB) *ExpenseService {
	return &ExpenseService{
		ExpenseDB: db,
	}
}

type ExpenseService struct {
	ExpenseDB db.DB
}

func (es *ExpenseService) GetExpense(id string) *model.Expense {
	for _, tempExpense := range es.ExpenseDB.Expenses {
		if tempExpense.Id == id {
			return &tempExpense
		}
	}
	return nil
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

func (es *ExpenseService) DeleteExpense(id string) error {

	for i, expense := range es.ExpenseDB.Expenses {
		if id == expense.Id {
			es.ExpenseDB.Expenses = append(es.ExpenseDB.Expenses[:i], es.ExpenseDB.Expenses[i+1:]...)
		}
	}
	return nil
}

// Stats Functions

func (es *ExpenseService) GetCurrentBalance(filter model.FilterFunc) float64 {
	sum := 0.0

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.EXPENSE {
			sum -= tempExpense.Amount
		} else if tempExpense.Type == model.INCOME {
			sum += tempExpense.Amount
		}

	}
	return sum
}

func (es *ExpenseService) GetCurrentExpenses(filter model.FilterFunc) float64 {
	sum := 0.0

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.EXPENSE {
			sum += tempExpense.Amount
		}

	}
	return sum
}

func (es *ExpenseService) GetCurrentIncomes(filter model.FilterFunc) float64 {
	sum := 0.0

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.INCOME {
			sum += tempExpense.Amount
		}

	}
	return sum
}

func (es *ExpenseService) GetLargestExpense(filter model.FilterFunc) model.Expense {
	var e model.Expense

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.INCOME {
			continue
		}
		if tempExpense.Amount > e.Amount {
			e = tempExpense
		}
	}
	return e
}

func (es *ExpenseService) GetExpensesByCategory(filter model.FilterFunc) ([]string, []float64) {
	categories := make(map[string]float64)

	var categoryNames []string
	var categoryValues []float64

	categoryNames = db.Categories

	for _, tempExpense := range es.ExpenseDB.Expenses {

		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

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

func (es *ExpenseService) GetExpensesByMonth(filter model.FilterFunc) ([]string, []float64) {

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
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.INCOME {
			continue
		}

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

func (es *ExpenseService) GetExpensesByDay(filter model.FilterFunc) map[string]float64 {
	days := make(map[string]float64)

	for _, tempExpense := range es.ExpenseDB.Expenses {

		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		day := tempExpense.Date.Day()
		month := tempExpense.Date.Month()
		year := tempExpense.Date.Year()

		key := fmt.Sprintf("%02d/%02d/%04d", day, month, year)

		days[key] += tempExpense.Amount
	}
	return days
}

func (es *ExpenseService) GetLongestStreakWithoutExpense(filter model.FilterFunc) int {

	var dates []time.Time

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		dates = append(dates, tempExpense.Date)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	maxStreak := 0

	if len(dates) == 1 {
		return 1
	}

	if len(dates) > 1 {
		firstDate := dates[0]
		for i := 1; i < len(dates); i++ {

			tempStreak := dates[i].Day() - firstDate.Day()

			if tempStreak > maxStreak {
				maxStreak = tempStreak
			}

			firstDate = dates[i]
		}
	}

	return maxStreak
}

func (es *ExpenseService) GetCountsByCategory(filter model.FilterFunc) ([]string, []float64) {

	categoryFrequencies := make(map[string]float64)

	var categoryNames []string
	var categoryCounts []float64

	categoryNames = db.Categories

	for _, tempExpense := range es.ExpenseDB.Expenses {
		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		categoryFrequencies[tempExpense.Category] += 1

		if !slices.Contains(categoryNames, tempExpense.Category) {
			categoryNames = append(categoryNames, tempExpense.Category)
		}
	}

	slices.Sort(categoryNames)

	for _, cName := range categoryNames {
		categoryCounts = append(categoryCounts, categoryFrequencies[cName])
	}

	return categoryNames, categoryCounts
}
