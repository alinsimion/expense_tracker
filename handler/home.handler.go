package handler

import (
	"sort"
	"strconv"

	"github.com/a-h/templ"
	"github.com/alinsimion/expense_tracker/model"
	expense "github.com/alinsimion/expense_tracker/view/expenses"
	home "github.com/alinsimion/expense_tracker/view/home"
	"github.com/labstack/echo/v4"
)

func (es *ExpenseHandler) ShowHomescreen(ctx echo.Context) error {

	var returnView templ.Component
	user := ctx.Request().Context().Value(userContextKey)
	if user != nil {
		if user.(*model.User).Email != "" {
			skip, err := strconv.Atoi(ctx.QueryParam("skip"))
			if err != nil {
				skip = 0
			}

			limit, err := strconv.Atoi(ctx.QueryParam("limit"))
			if err != nil {
				limit = 1000
			}

			expenses := es.store.GetExpenses(skip, limit, user.(*model.User).Id)

			sort.Slice(expenses, func(i, j int) bool {
				return expenses[i].Date.After(expenses[j].Date)
			})
			returnView = expense.ShowExpenseList("Expenses", expenses, []bool{false})
		}
	} else {
		returnView = home.ShowHomeWithLayout()
	}

	return View(ctx, returnView)
}
