package handler

import (
	"fmt"
	"net/http"
	"strconv"

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

	// return c.JSON(http.StatusOK, expense)
}

func (eh *ExpenseHandler) GetExpenses(c echo.Context) error {
	skip, _ := strconv.Atoi(c.QueryParam("skip"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	fmt.Println(skip)
	fmt.Println(limit)

	expenses := eh.ExpenseService.GetExpenses(skip, limit)

	se := view.ShowExpenseList("Expenses", view.ExpenseList(expenses))

	return eh.View(c, se)
}

func (eh *ExpenseHandler) AddExpense(c echo.Context) error {
	m := new(model.Expense)
	if err := c.Bind(m); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, m)
}

func (uh *ExpenseHandler) View(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
