package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/model"
	"github.com/alinsimion/expense_tracker/view"
	stats "github.com/alinsimion/expense_tracker/view/stats"
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

func (eh *ExpenseHandler) ShowStats(c echo.Context) error {

	start := c.QueryParams().Get("start")
	end := c.QueryParams().Get("end")
	if start == "" && end == "" {
		start = "01/01/" + fmt.Sprintf("%d", time.Now().Year())
		end = "31/12/" + fmt.Sprintf("%d", time.Now().Year())
	}
	// could get startDate and endDate for custom period stats
	user := c.Request().Context().Value(userContextKey).(*model.User)

	balance := eh.store.GetCurrentBalance(filter(start, end), user.Id)
	expenses := eh.store.GetCurrentExpenses(filter(start, end), user.Id)
	incomes := eh.store.GetCurrentIncomes(filter(start, end), user.Id)
	largestExpense := eh.store.GetLargestExpense(filter(start, end), user.Id)

	var categoryAmounts model.SortedMap
	categoryAmounts.Keys, categoryAmounts.Values = eh.store.GetExpensesByCategory(filter(start, end), user.Id)

	var monthMap model.SortedMap
	monthMap.Keys, monthMap.Values = eh.store.GetExpensesByMonth(filter(start, end), user.Id)

	days := eh.store.GetExpensesByDay(filter(start, end), user.Id)
	dayStreak := eh.store.GetLongestStreakWithoutExpense(filter(start, end), user.Id)

	var categoryCounts model.SortedMap
	categoryCounts.Keys, categoryCounts.Values = eh.store.GetCountsByCategory(filter(start, end), user.Id)

	ss := stats.ShowAllStats(balance, largestExpense, categoryAmounts, monthMap, dayStreak, days, categoryCounts, incomes, expenses)

	return View(c, ss)
}

func (eh *ExpenseHandler) ShowMonths(c echo.Context) error {
	yearParam := c.QueryParams().Get("year")
	user := c.Request().Context().Value(userContextKey).(*model.User)

	monthlyExpenses := make(map[string]model.SortedMap)
	monthlyIncomes := make(map[string]float64)

	monthNames := []string{"January", "February", "March", "April",
		"May", "June", "July", "August",
		"September", "October", "November", "December"}

	for index, monthName := range monthNames {
		var start string
		var end string

		year := time.Now().Year()

		if yearParam != "" {
			year, _ = strconv.Atoi(yearParam)
		}

		start = fmt.Sprintf("01/%0d/%d", index+1, year)
		end = fmt.Sprintf("31/%0d/%d", index+1, year)

		var monthMap model.SortedMap
		monthMap.Keys, monthMap.Values = eh.store.GetExpensesByCategory(filter(start, end), user.Id)
		expenseSum := 0.0

		for i := range monthMap.Values {
			expenseSum += monthMap.Values[i]
		}
		if expenseSum > 0 {
			monthlyExpenses[monthName] = monthMap
		}

		_, totals := eh.store.GetIncomeByCategory(filter(start, end), user.Id)

		// fmt.Println(year, " ", monthName, " ", totals)
		incomeSum := 0.0

		for i := range totals {
			incomeSum += totals[i]
		}

		monthlyIncomes[monthName] = incomeSum
		// fmt.Println(monthlyIncomes)
	}

	sm := view.ShowMonthsWithLayout(monthNames, monthlyExpenses, monthlyIncomes)

	return View(c, sm)
}
