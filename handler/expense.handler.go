package handler

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/view"
	"github.com/labstack/echo/v4"
)

type ExpenseServiceInterface interface {
	GetExpense(id string) model.Expense
	GetExpenses(skip int, limit int) []model.Expense
	AddExpense(e model.Expense)
	UpdateExpense(id string, e model.Expense) *model.Expense
	GetCurrentBalance() float64
	GetLargestExpense() model.Expense
	GetAmountByCategory() ([]string, []float64)
	GetAmountByMonth() ([]string, []float64)
	GetAmountByDay() map[string]float64
	GetLongestStreakWithoutExpense() int
}

func NewExpenseHandler(us ExpenseServiceInterface) *ExpenseHandler {
	return &ExpenseHandler{
		ExpenseService: us,
	}
}

type ExpenseHandler struct {
	ExpenseService ExpenseServiceInterface
}

func (eh *ExpenseHandler) ShowExpenseById(c echo.Context) error {
	expenseId := c.Param("id")

	fmt.Println(expenseId)

	expense := eh.ExpenseService.GetExpense(expenseId)

	se := view.ShowExpense(expense)

	return eh.View(c, se)
}

func (eh *ExpenseHandler) GetExpenses(c echo.Context) error {
	skip, err := strconv.Atoi(c.QueryParam("skip"))
	if err != nil {
		skip = 0
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 1000
	}

	expenses := eh.ExpenseService.GetExpenses(skip, limit)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	se := view.ShowExpenseList("Expenses", view.ExpenseList(expenses))

	return eh.View(c, se)
}

func (eh *ExpenseHandler) AddExpense(c echo.Context) error {

	fmt.Println(c.FormParams())

	m := new(model.Expense)
	var err error

	m.Description = c.FormValue("description")
	m.Amount, err = strconv.ParseFloat(c.FormValue("amount"), 64)

	if err != nil {
		return err
	}

	m.Category = c.FormValue("categories")
	m.Currency = "RON"
	m.Date = time.Now()
	m.Type = c.FormValue("type")

	fmt.Println(m)

	eh.ExpenseService.AddExpense(*m)

	return eh.ShowAddExpense(c)

	// return c.JSON(http.StatusCreated, m)
}

func (eh *ExpenseHandler) ShowStats(c echo.Context) error {

	// could get startDate and endDate for custom period stats

	balance := eh.ExpenseService.GetCurrentBalance()
	largestExpense := eh.ExpenseService.GetLargestExpense()

	var categoryMap model.SortedMap

	categoryMap.Keys, categoryMap.Values = eh.ExpenseService.GetAmountByCategory()

	var monthMap model.SortedMap

	monthMap.Keys, monthMap.Values = eh.ExpenseService.GetAmountByMonth()
	days := eh.ExpenseService.GetAmountByDay()
	dayStreak := eh.ExpenseService.GetLongestStreakWithoutExpense()

	ss := view.ShowAllStats(balance, largestExpense, categoryMap, monthMap, dayStreak, days)

	return eh.View(c, ss)

}

func (eh *ExpenseHandler) ShowAddExpense(c echo.Context) error {

	sae := view.Page()

	return eh.View(c, sae)

}

func (uh *ExpenseHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
