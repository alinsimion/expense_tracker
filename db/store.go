package db

import "github.com/alinsimion/expense_tracker/model"

type Store interface {
	GetExpenseById(id string) *model.Expense
	GetExpenses(skip int, limit int, uid int64) []model.Expense
	CreateExpense(e model.Expense)
	UpdateExpense(id string, e model.Expense) (int64, error)
	DeleteExpense(id string) error

	GetCurrentBalance(filter model.FilterFunc, userId int64) float64
	GetLargestExpense(filter model.FilterFunc, userId int64) model.Expense
	GetExpensesByCategory(filter model.FilterFunc, userId int64) ([]string, []float64)
	GetExpensesByMonth(filter model.FilterFunc, userId int64) ([]string, []float64)
	GetExpensesByDay(filter model.FilterFunc, userId int64) map[string]float64
	GetLongestStreakWithoutExpense(filter model.FilterFunc, userId int64) int
	GetCountsByCategory(filter model.FilterFunc, userId int64) ([]string, []float64)
	GetCurrentIncomes(filter model.FilterFunc, userId int64) float64
	GetCurrentExpenses(filter model.FilterFunc, userId int64) float64
	GetCategories(userId int64) []string

	GetUserById(id string) *model.User
	GetUserByEmail(email string) *model.User
	GetUsers(skip int, limit int) []model.User
	CreateUser(e model.User)
	UpdateUser(id string, e model.User) (int64, error)
}
