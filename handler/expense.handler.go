package handler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/service"
	"github.com/alinsimion/expense_tracker/view"
	"github.com/labstack/echo/v4"
)

type ExpenseHandler struct {
	ExpenseService service.ExpenseServiceInterface
}

func NewExpenseHandler(us service.ExpenseServiceInterface) *ExpenseHandler {
	return &ExpenseHandler{
		ExpenseService: us,
	}
}

// View Functions

func (eh *ExpenseHandler) GetExpenseById(c echo.Context) error {
	expenseId := c.Param("id")

	fmt.Println(expenseId)

	expense := *eh.ExpenseService.GetExpense(expenseId)

	se := view.ShowExpense(expense)

	return View(c, se)
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

	return View(c, se)
}

func (eh *ExpenseHandler) AddExpense(c echo.Context) error {

	var income string

	if c.FormValue("type") == "" {
		income = model.EXPENSE
	} else {
		income = model.INCOME
	}

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
	m.Type = income

	eh.ExpenseService.AddExpense(*m)
	// return c.JSON(http.StatusCreated, m)
	return eh.ShowAddExpense(c)
}

func (eh *ExpenseHandler) DeleteExpense(c echo.Context) error {
	return nil
}

// API Functions
// TODO

func filter(start string, end string) model.FilterFunc {
	return func(e model.Expense) bool {

		var startDate time.Time
		var endDate time.Time

		if start != "" && end != "" {

			d, _ := strconv.Atoi(strings.Split(start, "/")[0])
			m, _ := strconv.Atoi(strings.Split(start, "/")[1])
			y, _ := strconv.Atoi(strings.Split(start, "/")[2])

			startDate = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

			d, _ = strconv.Atoi(strings.Split(end, "/")[0])
			m, _ = strconv.Atoi(strings.Split(end, "/")[1])
			y, _ = strconv.Atoi(strings.Split(end, "/")[2])

			endDate = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
		}

		if startDate.IsZero() || endDate.IsZero() {
			return false
		}

		if (e.Date.After(startDate) || e.Date.Equal(startDate)) && (e.Date.Before(endDate) || e.Date.Equal(endDate)) {
			return false
		}

		return true
	}
}

func (eh *ExpenseHandler) ShowAddExpense(c echo.Context) error {
	sae := view.Page()

	return View(c, sae)
}

func (eh *ExpenseHandler) ShowStats(c echo.Context) error {

	start := c.QueryParams().Get("start")
	end := c.QueryParams().Get("end")
	// could get startDate and endDate for custom period stats

	balance := eh.ExpenseService.GetCurrentBalance(filter(start, end))
	expenses := eh.ExpenseService.GetCurrentExpenses(filter(start, end))
	incomes := eh.ExpenseService.GetCurrentIncomes(filter(start, end))
	largestExpense := eh.ExpenseService.GetLargestExpense(filter(start, end))

	var categoryAmounts model.SortedMap
	categoryAmounts.Keys, categoryAmounts.Values = eh.ExpenseService.GetAmountByCategory(filter(start, end))

	var monthMap model.SortedMap
	monthMap.Keys, monthMap.Values = eh.ExpenseService.GetAmountByMonth(filter(start, end))

	days := eh.ExpenseService.GetAmountByDay(filter(start, end))
	dayStreak := eh.ExpenseService.GetLongestStreakWithoutExpense(filter(start, end))

	var categoryCounts model.SortedMap
	categoryCounts.Keys, categoryCounts.Values = eh.ExpenseService.GetCountsByCategory(filter(start, end))

	ss := view.ShowAllStats(balance, largestExpense, categoryAmounts, monthMap, dayStreak, days, categoryCounts, incomes, expenses)

	return View(c, ss)
}

func View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
