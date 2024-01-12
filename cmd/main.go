package main

import (
	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/handler"
	"github.com/alinsimion/expense_tracker/service"
	"github.com/labstack/echo/v4"
)

func main() {
	db := db.OpenDB()

	e := echo.New()

	expenseService := service.NewExpenseService(db)

	expenseHandler := handler.NewExpenseHandler(expenseService)

	handler.SetupRoutes(e, expenseHandler)

	e.Logger.Fatal(e.Start(":1323"))

}
