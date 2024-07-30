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

// func (eh *ExpenseHandler) ShowExpenseById(c echo.Context) error {
// 	expenseId := c.Param("id")

// 	// fmt.Println(expenseId)

// 	expense := *eh.store.GetExpenseById(expenseId)

// 	tempList := []model.Expense{expense}

// 	se := expenseslib.ShowExpenseList("", tempList, []bool{false}, )

// 	return View(c, se)
// }

func (eh *ExpenseHandler) ShowExpenses(c echo.Context) error {

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

	var pages []bool

	for i := 1; i <= len(expenses); i += pageLimit {
		pages = append(pages, false)
	}
	pages[page] = true

	categories := eh.store.GetCategories(user.Id)

	se := expenseslib.ShowExpenseList("Expenses", filteredExpenseList, pages, categories)

	return View(c, se)
}

type ExpenseFilter struct {
	SearchTerm string
	Category   string
	Type       string
	StartDate  time.Time
	EndDate    time.Time
	MinAmount  float64
	MaxAmount  float64
}

func (ef *ExpenseFilter) filterExpenses(expenses []model.Expense) []model.Expense {

	var tempExpenses []model.Expense
	var tempMatchCount int
	var totalMatchCount int

	for _, expense := range expenses {
		if ef.SearchTerm != "" {
			totalMatchCount += 1
			str := strings.ToLower(expense.Description)
			substr := strings.ToLower(ef.SearchTerm)
			if strings.Contains(str, substr) {
				tempMatchCount += 1
			}
		}

		if ef.Category != "Pick a category" && ef.Category != "All" {
			totalMatchCount += 1
			if ef.Category == expense.Category {
				tempMatchCount += 1
			}
		}

		if !ef.StartDate.IsZero() {
			totalMatchCount += 1
			if ef.StartDate.Before(expense.Date) {
				tempMatchCount += 1
			}
		}

		if !ef.EndDate.IsZero() {
			totalMatchCount += 1
			if ef.EndDate.After(expense.Date) {
				tempMatchCount += 1
			}
		}

		if ef.Type != "Pick a type" {
			totalMatchCount += 1

			if ef.Type == expense.Type || ef.Type == model.BOTH {
				tempMatchCount += 1
			}
		}

		if ef.MinAmount != 0 {
			totalMatchCount += 1

			if ef.MinAmount <= expense.Amount {
				tempMatchCount += 1
			}
		}

		if ef.MaxAmount != 0 {
			totalMatchCount += 1

			if ef.MaxAmount > expense.Amount {
				tempMatchCount += 1
			}
		}

		if tempMatchCount == totalMatchCount {
			tempExpenses = append(tempExpenses, expense)
		}

		tempMatchCount = 0
		totalMatchCount = 0

	}

	return tempExpenses
}

func (eh *ExpenseHandler) ShowExpensesTable(c echo.Context) error {

	var onePageExpenses []model.Expense
	var totalAmount float64

	var startDateTime time.Time
	var endDateTime time.Time

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

	c.Request().ParseForm()
	search := c.Request().Form.Get("search")
	category := c.Request().Form.Get("category")
	expenseType := c.Request().Form.Get("expense_type")

	startDate := c.Request().Form.Get("start_date")
	endDate := c.Request().Form.Get("end_date")

	minAmount, _ := strconv.ParseFloat(c.Request().Form.Get("min_amount"), 64)
	maxAmount, _ := strconv.ParseFloat(c.Request().Form.Get("max_amount"), 64)

	if startDate != "" {
		y, _ := strconv.Atoi(strings.Split(startDate, "-")[0])
		m, _ := strconv.Atoi(strings.Split(startDate, "-")[1])
		d, _ := strconv.Atoi(strings.Split(startDate, "-")[2])

		startDateTime = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)

	}

	if endDate != "" {
		y, _ := strconv.Atoi(strings.Split(endDate, "-")[0])
		m, _ := strconv.Atoi(strings.Split(endDate, "-")[1])
		d, _ := strconv.Atoi(strings.Split(endDate, "-")[2])

		endDateTime = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	}

	user := c.Request().Context().Value(userContextKey).(*model.User)

	expenses := eh.store.GetExpenses(0, 0, user.Id)

	ef := ExpenseFilter{
		SearchTerm: search,
		Category:   category,
		Type:       expenseType,
		StartDate:  startDateTime,
		EndDate:    endDateTime,
		MinAmount:  minAmount,
		MaxAmount:  maxAmount,
	}

	expenses = ef.filterExpenses(expenses)

	for _, e := range expenses {
		totalAmount += e.Amount
	}

	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date.After(expenses[j].Date)
	})

	if nextFor != 0 {
		for idx, expense := range expenses {
			if expense.Id == int64(nextFor) {
				start := min(len(expenses)-pageLimit, idx+1)
				end := min(start+pageLimit, len(expenses))

				onePageExpenses = expenses[start:end]
				page = start / pageLimit

				break
			}
		}
	} else if prevFor != 0 {
		for idx, expense := range expenses {
			if expense.Id == int64(prevFor) {
				start := max(0, idx-pageLimit)
				end := min(start+pageLimit, len(expenses))
				onePageExpenses = expenses[start:end]

				page = start / pageLimit

				break
			}
		}
	} else {
		onePageExpenses = expenses[page*pageLimit : min((page+1)*pageLimit, len(expenses))]
	}

	var pages []bool

	for i := 1; i <= len(expenses); i += pageLimit {
		pages = append(pages, false)
	}
	if len(pages) > 0 {
		pages[page] = true
	}

	se := expenseslib.ShowExpenseTableBody(onePageExpenses, pages, totalAmount, len(expenses))

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

	categories := eh.store.GetCategories(user.Id)
	return View(c, expenseslib.ExpenseList2(expenses, []bool{false}, categories))
}
