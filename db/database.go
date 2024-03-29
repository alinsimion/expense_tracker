package db

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/model"
)

type DB struct {
	Expenses []model.Expense
}

var db DB

var Categories []string

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

	for id, eachrecord := range records {
		amount, _ := strconv.ParseFloat(eachrecord[3], 64)

		date := strings.Split(eachrecord[1], ",")

		d, _ := strconv.Atoi(strings.Split(date[1], "/")[0])
		m, _ := strconv.Atoi(strings.Split(date[1], "/")[1])
		y, _ := strconv.Atoi(strings.Split(date[1], "/")[2])

		tempTime := time.Date(2000+y, time.Month(m), d, 0, 0, 0, 0, time.Local)

		// category := strings.Replace(eachrecord[2], " ", "", -1)
		category := eachrecord[2]
		if !slices.Contains(Categories, category) {
			Categories = append(Categories, category)
		}

		tempExpense := model.Expense{
			Id:          strconv.Itoa(id),
			Description: eachrecord[0],
			Amount:      amount,
			Category:    category,
			Date:        tempTime,
			Currency:    "RON",
			Type:        model.EXPENSE,
		}
		expenses = append(expenses, tempExpense)
	}

	var file2 *os.File
	file2, err = os.Open("/Users/asimion/Desktop/Personal/Projects/expense_tracker/test_data/temp_income.csv")

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	defer file.Close()

	reader2 := csv.NewReader(file2)

	var records2 [][]string
	records2, err = reader2.ReadAll()

	if err != nil {
		fmt.Println("Error reading records")
	}

	for id, eachrecord := range records2 {
		amount, _ := strconv.ParseFloat(eachrecord[3], 64)

		date := strings.Split(eachrecord[1], ",")

		d, _ := strconv.Atoi(strings.Split(date[1], "/")[0])
		m, _ := strconv.Atoi(strings.Split(date[1], "/")[1])
		y, _ := strconv.Atoi(strings.Split(date[1], "/")[2])

		tempTime := time.Date(2000+y, time.Month(m), d, 0, 0, 0, 0, time.Local)

		// category := strings.Replace(eachrecord[2], " ", "", -1)
		category := eachrecord[2]
		if !slices.Contains(Categories, category) {
			Categories = append(Categories, category)
		}

		tempExpense := model.Expense{
			Id:          strconv.Itoa(id),
			Description: eachrecord[0],
			Amount:      amount,
			Category:    category,
			Date:        tempTime,
			Currency:    "RON",
			Type:        model.INCOME,
		}
		expenses = append(expenses, tempExpense)
	}

	db = DB{
		Expenses: expenses,
	}

	fmt.Println("Am deschis DB")
	return db
}
