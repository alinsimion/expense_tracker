package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/util"
	"github.com/alinsimion/expense_tracker/view"
	expenseslib "github.com/alinsimion/expense_tracker/view/expenses"
	"github.com/labstack/echo/v4"
)

const (
	pageLimit = 30
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

	// fmt.Println(expenseId)

	expense := *eh.store.GetExpenseById(expenseId)

	tempList := []model.Expense{expense}

	se := expenseslib.ShowExpenseList("", tempList, []bool{false})

	return View(c, se)
}

func (eh *ExpenseHandler) ShowExpenses(c echo.Context) error {

	// skip, err := strconv.Atoi(c.QueryParam("skip"))
	// if err != nil {
	// 	skip = 0
	// }

	// limit, err := strconv.Atoi(c.QueryParam("limit"))
	// if err != nil {
	// 	limit = 1000
	// }

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	user := c.Request().Context().Value(userContextKey).(*model.User)

	// Ar trebui adaugat si filtru
	expenses := eh.store.GetExpenses(0, 0, user.Id)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	filteredExpenseList := expenses[page*pageLimit : (page+1)*pageLimit]

	// sort.Slice(filteredExpenseList, func(i, j int) bool {
	// 	return filteredExpenseList[i].Date.After(filteredExpenseList[j].Date)
	// })

	var pages []bool

	for i := 1; i <= len(expenses); i += pageLimit {
		pages = append(pages, false)
	}
	pages[page] = true

	se := expenseslib.ShowExpenseList("Expenses", filteredExpenseList, pages)

	return View(c, se)
}

func (eh *ExpenseHandler) ShowExpensesTable(c echo.Context) error {

	var filteredExpenseList []model.Expense
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 0
	}

	nextFor, err := strconv.Atoi(c.QueryParam("next_for"))
	if err != nil {
		nextFor = 0
	}

	prevFor, err := strconv.Atoi(c.QueryParam("prev_for"))
	if err != nil {
		prevFor = 0
	}

	user := c.Request().Context().Value(userContextKey).(*model.User)

	// Ar trebui adaugat si filtru
	expenses := eh.store.GetExpenses(0, 0, user.Id)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	if nextFor != 0 {
		for idx, expense := range expenses {
			if expense.Id == int64(nextFor) {
				start := min(len(expenses)-pageLimit, idx+1)
				end := min(start+pageLimit, len(expenses))
				filteredExpenseList = expenses[start:end]

				page = (len(expenses) - start) / pageLimit

				break
			}
		}
	} else if prevFor != 0 {
		for idx, expense := range expenses {
			if expense.Id == int64(prevFor) {
				start := max(0, idx-pageLimit)
				end := min(start+pageLimit, len(expenses))
				filteredExpenseList = expenses[start:end]

				page = (len(expenses) - start) / pageLimit

				break
			}
		}
	} else {
		filteredExpenseList = expenses[page*pageLimit : (page+1)*pageLimit]
	}

	// sort.Slice(filteredExpenseList, func(i, j int) bool {
	// 	return filteredExpenseList[i].Date.After(filteredExpenseList[j].Date)
	// })

	var pages []bool

	for i := 1; i <= len(expenses); i += pageLimit {
		pages = append(pages, false)
	}
	pages[page] = true

	se := expenseslib.ShowExpenseTableBody(filteredExpenseList, pages)

	return View(c, se)
}

// Danger
func (eh *ExpenseHandler) UnsecureAddExpense(c echo.Context) error {
	var addParams view.AddExpenseParams

	err := c.Bind(&addParams)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(addParams)

	tempExpense := model.Expense{
		Uid:         1,
		Description: addParams.Description,
		Category:    addParams.Category,
		Amount:      addParams.Amount,
		Date:        addParams.Date,
		Currency:    addParams.Currency,
		Type:        addParams.ExpenseType,
	}

	eh.store.CreateExpense(tempExpense)

	return nil
}

func (eh *ExpenseHandler) ShowAddExpense(c echo.Context) error {
	if c.Request().Method == "POST" {

		var addParamas view.AddExpenseParams

		if c.FormValue("type") == "" {
			addParamas.ExpenseType = model.EXPENSE
		} else {
			addParamas.ExpenseType = model.INCOME
		}

		// e := new(model.Expense)
		var err error

		addParamas.Description = c.FormValue("description")
		addParamas.Amount, err = strconv.ParseFloat(c.FormValue("amount"), 64)

		if err != nil {

			return err
		}
		newCategory := c.FormValue("category")

		if newCategory != "" {
			addParamas.Category = newCategory
		} else {
			addParamas.Category = c.FormValue("categories")
		}
		addParamas.Currency = c.FormValue("currency")

		if strings.Contains(c.FormValue("date"), "/") {
			d, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[0])
			m, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[1])
			y, _ := strconv.Atoi(strings.Split(c.FormValue("date"), "/")[2])

			tempTime := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

			addParamas.Date = tempTime
		} else {
			addParamas.Date = time.Now()
		}

		// addParamas.Type = income

		tempExpense := model.Expense{
			Uid:         int64(c.Request().Context().Value(userContextKey).(*model.User).Id),
			Description: addParamas.Description,
			Category:    addParamas.Category,
			Amount:      addParamas.Amount,
			Date:        addParamas.Date,
			Currency:    addParamas.Currency,
			Type:        addParamas.ExpenseType,
		}

		if c.FormValue("edit_id") != "0" {

			_, err := eh.store.UpdateExpense(c.FormValue("edit_id"), tempExpense)

			if err != nil {
				slog.Error("Error while updating expense", "err", err.Error())
			}
		} else {
			eh.store.CreateExpense(tempExpense)

		}

		http.Redirect(c.Response(), c.Request(), util.GetFullUrl("/expense"), http.StatusSeeOther)
		return nil
	}

	return eh.ShowAddExpenseWithLayout(c)
}

func (eh *ExpenseHandler) ShowAddExpenseWithLayout(c echo.Context) error {
	user := c.Request().Context().Value(userContextKey).(*model.User)

	sae := view.ShowAddExpenseForm(eh.store.GetCategories(user.Id), view.AddExpenseParams{
		Date: time.Now(),
	}, view.AddExpenseErrors{})

	return View(c, sae)
}

func (eh *ExpenseHandler) ShowEditExpenseById(c echo.Context) error {
	user := c.Request().Context().Value(userContextKey).(*model.User)

	expenseId := c.Param("id")

	expense := eh.store.GetExpenseById(expenseId)

	addParamas := view.AddExpenseParams{
		Id:          expense.Id,
		Currency:    expense.Currency,
		Category:    expense.Category,
		Description: expense.Description,
		Amount:      expense.Amount,
		Date:        expense.Date,
		ExpenseType: expense.Type,
	}

	// fmt.Println(addParamas)

	sae := view.ShowAddExpenseForm(eh.store.GetCategories(user.Id), addParamas, view.AddExpenseErrors{})

	return View(c, sae)
}

func (eh *ExpenseHandler) ShowDeleteExpense(c echo.Context) error {
	expenseId := c.Param("id")

	// check if expense with id belongs to auth user - just paranoid

	eh.store.DeleteExpense(expenseId)

	user := c.Request().Context().Value(userContextKey).(*model.User)

	expenses := eh.store.GetExpenses(0, 0, user.Id)

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	var pages []bool

	for i := 1; i <= len(expenses); i += pageLimit {
		pages = append(pages, false)
	}
	pages[0] = true

	return View(c, expenseslib.ExpenseList2(expenses, []bool{false}))
}
