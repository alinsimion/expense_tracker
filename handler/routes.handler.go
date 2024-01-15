package handler

import "github.com/labstack/echo/v4"

func SetupRoutes(app *echo.Echo, h *ExpenseHandler) {
	group := app.Group("/expense")

	group.GET("/:id", h.ShowExpenseById)
	group.GET("", h.GetExpenses)
	group.POST("", h.AddExpense)

	group2 := app.Group("/stats")

	group2.GET("", h.ShowStats)

	group3 := app.Group("/add_expense")

	group3.GET("", h.ShowAddExpense)
	group3.POST("", h.AddExpense)
}
