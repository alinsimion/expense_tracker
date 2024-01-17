package handler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

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

func (eh *ExpenseHandler) ShowExpenseById(c echo.Context) error {
	expenseId := c.Param("id")

	fmt.Println(expenseId)

	expense := *eh.ExpenseService.GetExpense(expenseId)

	se := view.ShowExpense(expense)

	return View(c, se)
}

func (eh *ExpenseHandler) ShowExpenses(c echo.Context) error {
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

func (eh *ExpenseHandler) ShowAddExpense(c echo.Context) error {

	if c.Request().Method == "POST" {
		var income string

		if c.FormValue("type") == "" {
			income = model.EXPENSE
		} else {
			income = model.INCOME
		}

		e := new(model.Expense)
		var err error

		e.Description = c.FormValue("description")
		e.Amount, err = strconv.ParseFloat(c.FormValue("amount"), 64)

		if err != nil {
			return err
		}

		e.Category = c.FormValue("categories")
		e.Currency = "RON"

		if strings.Contains(c.FormValue("date"), "/") {
			d, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[0])
			m, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[1])
			y, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[2])

			tempTime := time.Date(2000+y, time.Month(m), d, 0, 0, 0, 0, time.Local)

			e.Date = tempTime
		} else {
			e.Date = time.Now()
		}

		e.Type = income

		eh.ExpenseService.AddExpense(*e)
		// return c.JSON(http.StatusCreated, m)
	}

	return eh.ShowAddExpenseWithLayout(c)
}

func (eh *ExpenseHandler) ShowEditExpenseById(c echo.Context) error {
	expenseId := c.Param("id")

	expense := eh.ExpenseService.GetExpense(expenseId)

	eh.ExpenseService.AddExpense(*expense)

	return eh.ShowAddExpense(c)
}

func (eh *ExpenseHandler) ShowDeleteExpense(c echo.Context) error {
	expenseId := c.Param("id")

	eh.ExpenseService.DeleteExpense(expenseId)

	expenses := eh.ExpenseService.GetExpenses(0, 500)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	return View(c, view.ExpenseList(expenses))
}

// API Functions
// TODO

func (eh *ExpenseHandler) ApiGetExpenseById(c echo.Context) error {
	return nil
}

func (eh *ExpenseHandler) ApiGetExpensesWithFilter(c echo.Context) error {
	return nil
}

func (eh *ExpenseHandler) ApiAddExpense(c echo.Context) error {
	return nil
}

func (eh *ExpenseHandler) ApiDeleteExpense(c echo.Context) error {
	return nil
}

// Patch ?
func (eh *ExpenseHandler) ApiEditExpense(c echo.Context) error {
	return nil
}
