package db

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/model"
)

type DB struct {
	Expenses []model.Expense
}

var db DB

func OpenDB() DB {

	var expenses []model.Expense

	file, err := os.Open("/Users/asimion/Desktop/Personal/Projects/expense_tracker/test_data/temp_data.csv")

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading records")
	}

	for _, eachrecord := range records {
		amount, _ := strconv.Atoi(eachrecord[3])

		date := strings.Split(eachrecord[1], ",")

		d, _ := strconv.Atoi(strings.Split(date[1], "/")[0])
		m, _ := strconv.Atoi(strings.Split(date[1], "/")[1])
		y, _ := strconv.Atoi(strings.Split(date[1], "/")[2])

		tempTime := time.Date(2000+y, time.Month(m), d, 0, 0, 0, 0, time.Local)

		tempExpense := model.Expense{
			Description: eachrecord[0],
			Amount:      amount,
			Category:    eachrecord[2],
			Date:        tempTime,
			Currency:    "RON",
		}
		expenses = append(expenses, tempExpense)
	}

	db = DB{
		Expenses: expenses,
	}

	fmt.Println("Am deschis DB")
	return db
}