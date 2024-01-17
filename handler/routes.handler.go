package handler

import "github.com/labstack/echo/v4"

func SetupRoutes(app *echo.Echo, h *ExpenseHandler) {
	// View
	group := app.Group("/expense")
	group.GET("", h.ShowExpenses)
	group.GET("/:id/details", h.ShowExpenseById)
	group.DELETE("/:id/delete", h.ShowDeleteExpense)
	group.GET("/add", h.ShowAddExpense)
	group.POST("/add", h.ShowAddExpense)
	group.GET("/:id/edit", h.ShowEditExpenseById)
	// group.GET("/:id/details", h.ShowEditExpenseById)

	// API
	group2 := app.Group("/api/")
	group2.GET("/", h.ApiGetExpensesWithFilter)
	group2.GET("/:id", h.ApiGetExpenseById)
	group2.DELETE("/:id", h.ApiDeleteExpense)
	group2.POST("/", h.ApiAddExpense)
	group2.PATCH("/:id", h.ApiEditExpense)

	// group.POST("", h.ShowAddExpense)

	group3 := app.Group("/stats")
	group3.GET("", h.ShowStats)

	// group4 := app.Group("/add_expense")

	// group4.GET("/:id", h.ShowAddExpense)
	// group4.POST("", h.ShowAddExpense)

}
