package handler

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo, h *ExpenseHandler) {
	// View

	app.GET("/signup", h.HandleSignupIndex)
	app.POST("/signup", h.HandleSignupCreate)
	app.GET("/login", h.HandleLoginIndex)
	app.POST("/login", h.HandleLoginCreate)
	app.GET("/login/provider/google", h.HandleLoginCreateProvider)
	app.POST("/logout", h.HandleLogoutCreate)
	app.GET("/auth/callback", h.HandleAuthCallback)

	app.GET("/", h.ShowHomescreen, h.WithUser)

	// app.POST("/add_expense", h.UnsecureAddExpense)

	group := app.Group("/expense", h.WithAuth)

	group.GET("", h.ShowExpenses)
	group.GET("/body", h.ShowExpensesTable)
	group.GET("/add", h.ShowAddExpense)
	group.POST("/add", h.ShowAddExpense)
	group.GET("/:id/edit", h.ShowEditExpenseById)
	group.GET("/:id/details", h.ShowExpenseById)
	group.DELETE("/:id/delete", h.ShowDeleteExpense)

	group3 := app.Group("/stats", h.WithAuth)
	group3.GET("", h.ShowStats)

	group4 := app.Group("/months", h.WithAuth)
	group4.GET("", h.ShowMonths)

}
