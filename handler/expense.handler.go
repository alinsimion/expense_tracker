package handler

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/util"
	expenseslib "github.com/alinsimion/expense_tracker/view/expenses"
	"github.com/labstack/echo/v4"
)

type ExpenseHandler struct {
	store db.Store
}

func NewExpenseServer(store db.Store) *ExpenseHandler {
	return &ExpenseHandler{
		store: store,
	}
}

func (eh *ExpenseHandler) ShowExpenseById(c echo.Context) error {
	expenseId := c.Param("id")

	expense := *eh.store.GetExpenseById(expenseId)

	tempList := []model.Expense{expense}

	se := expenseslib.ShowExpenseList("", tempList)

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

	user := c.Request().Context().Value(userContextKey).(*model.User)

	expenses := eh.store.GetExpenses(skip, limit, user.Id)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	se := expenseslib.ShowExpenseList("Expenses", expenses)

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
		newCategory := c.FormValue("category")

		if newCategory != "" {
			e.Category = newCategory
		} else {
			e.Category = c.FormValue("categories")
		}
		e.Currency = c.FormValue("currency")

		if strings.Contains(c.FormValue("date"), "/") {
			d, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[0])
			m, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[1])
			y, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[2])

			tempTime := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

			e.Date = tempTime
		} else {
			e.Date = time.Now()
		}

		e.Type = income

		e.Uid = int64(c.Request().Context().Value(userContextKey).(*model.User).Id)

		eh.store.CreateExpense(*e)
		http.Redirect(c.Response(), c.Request(), util.GetFullUrl("/expense"), http.StatusSeeOther)
		return nil
	}

	return eh.ShowAddExpenseWithLayout(c)

}

func (eh *ExpenseHandler) ShowEditExpenseById(c echo.Context) error {
	expenseId := c.Param("id")

	expense := eh.store.GetExpenseById(expenseId)

	eh.store.UpdateExpense(expenseId, *expense)

	return eh.ShowAddExpense(c)
}

func (eh *ExpenseHandler) ShowDeleteExpense(c echo.Context) error {
	expenseId := c.Param("id")

	eh.store.DeleteExpense(expenseId)

	user := c.Request().Context().Value(userContextKey).(*model.User)

	expenses := eh.store.GetExpenses(0, 500, user.Id)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	return View(c, expenseslib.ExpenseList2(expenses))
}
