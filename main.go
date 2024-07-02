package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/handler"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
	"github.com/nedpals/supabase-go"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load env variables")
	}

	HTTP_PORT := "1329"

	sbHost := os.Getenv("SUPABASE_URL")
	sbSecret := os.Getenv("SUPABASE_SECRET")

	handler.Client = supabase.CreateClient(sbHost, sbSecret)

	e := echo.New()
	e.Static("", "./css")

	// store := db.NewFileStore("temp_income.csv")

	store := db.NewSqlLiteStore("")

	expenseServer := handler.NewExpenseServer(store)

	handler.SetupRoutes(e, expenseServer)

	gothic.Store = sessions.NewCookieStore([]byte("asadasdas"))

	e.Use(session.Middleware(gothic.Store))

	// var callbackUrl string/

	// if os.Getenv("ENVIRONMENT") == "PRD" {
	// 	callbackUrl = "https://zportofolio.dev/etracker/auth/google/callback/"
	// } else if os.Getenv("ENVIRONMENT") == "DEV" {
	// 	callbackUrl = "http://localhost:1329/auth/google/callback"
	// }

	// goth.UseProviders(
	// 	google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), callbackUrl),
	// )

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", HTTP_PORT)))
}
