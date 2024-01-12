package handler

import "github.com/labstack/echo/v4"

func SetupRoutes(app *echo.Echo, h *ExpenseHandler) {
	group := app.Group("/expense")

	group.GET("/:id", h.ShowExpenseById)
	group.GET("", h.GetExpenses)
	group.POST("", h.AddExpense)
}
