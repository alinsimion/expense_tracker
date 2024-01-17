package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/view"
	"github.com/labstack/echo/v4"
)

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

func (eh *ExpenseHandler) ShowAddExpenseWithLayout(c echo.Context) error {
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
	categoryAmounts.Keys, categoryAmounts.Values = eh.ExpenseService.GetExpensesByCategory(filter(start, end))

	var monthMap model.SortedMap
	monthMap.Keys, monthMap.Values = eh.ExpenseService.GetExpensesByMonth(filter(start, end))

	days := eh.ExpenseService.GetExpensesByDay(filter(start, end))
	dayStreak := eh.ExpenseService.GetLongestStreakWithoutExpense(filter(start, end))

	var categoryCounts model.SortedMap
	categoryCounts.Keys, categoryCounts.Values = eh.ExpenseService.GetCountsByCategory(filter(start, end))

	ss := view.ShowAllStats(balance, largestExpense, categoryAmounts, monthMap, dayStreak, days, categoryCounts, incomes, expenses)

	return View(c, ss)
}
