package db

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alinsimion/expense_tracker/model"
)

type FileStore struct {
	fileName   string
	expenses   []model.Expense
	categories []string
}

func NewFileStore(fileName string) *FileStore {

	fileStore := FileStore{
		fileName:   fileName,
		expenses:   []model.Expense{},
		categories: []string{},
	}

	file, err := os.Open("test_data/temp_data.csv")

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
		if !slices.Contains(fileStore.categories, category) {
			fileStore.categories = append(fileStore.categories, category)
		}

		tempExpense := model.Expense{
			Id:          int64(id),
			Description: eachrecord[0],
			Amount:      amount,
			Category:    category,
			Date:        tempTime,
			Currency:    "RON",
			Type:        model.EXPENSE,
		}
		fileStore.expenses = append(fileStore.expenses, tempExpense)
	}

	var file2 *os.File
	file2, err = os.Open("test_data/temp_income.csv")

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
		if !slices.Contains(fileStore.categories, category) {
			fileStore.categories = append(fileStore.categories, category)
		}

		tempExpense := model.Expense{
			Id:          int64(id),
			Description: eachrecord[0],
			Amount:      amount,
			Category:    category,
			Date:        tempTime,
			Currency:    "RON",
			Type:        model.INCOME,
		}
		fileStore.expenses = append(fileStore.expenses, tempExpense)
	}

	// fmt.Println(fileStore.categories)

	return &fileStore

}

func (es *FileStore) GetExpense(id string) *model.Expense {
	for _, tempExpense := range es.expenses {
		number, _ := strconv.Atoi(id)
		if tempExpense.Id == int64(number) {
			return &tempExpense
		}
	}
	return nil
}

func (es *FileStore) GetExpenses(skip int, limit int) []model.Expense {
	maxLength := len(es.expenses)

	end := skip + limit

	if end > maxLength {
		end = maxLength
	}

	return es.expenses[skip:end]
}

func (es *FileStore) CreateExpense(e model.Expense) {
	es.expenses = append(es.expenses, e)
}

func (es *FileStore) UpdateExpense(id string, e model.Expense) (int64, error) {
	for i := 0; i < len(es.expenses); i++ {
		if es.expenses[i].Description == id {
			es.expenses[i] = e
			return 1, nil
		}
	}
	return 0, fmt.Errorf("Could not update expense")
}

func (es *FileStore) DeleteExpense(id string) error {

	for i, expense := range es.expenses {
		number, _ := strconv.Atoi(id)
		if expense.Id == int64(number) {
			es.expenses = append(es.expenses[:i], es.expenses[i+1:]...)
		}
	}
	return nil
}

// Stats Functions

func (es *FileStore) GetCurrentBalance(filter model.FilterFunc) float64 {
	sum := 0.0

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.EXPENSE {
			sum -= tempExpense.Amount
		} else if tempExpense.Type == model.INCOME {
			sum += tempExpense.Amount
		}

	}
	return sum
}

func (es *FileStore) GetCurrentExpenses(filter model.FilterFunc) float64 {
	sum := 0.0

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.EXPENSE {
			sum += tempExpense.Amount
		}

	}
	return sum
}

func (es *FileStore) GetCurrentIncomes(filter model.FilterFunc) float64 {
	sum := 0.0

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.INCOME {
			sum += tempExpense.Amount
		}

	}
	return sum
}

func (es *FileStore) GetLargestExpense(filter model.FilterFunc) model.Expense {
	var e model.Expense

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.INCOME {
			continue
		}
		if tempExpense.Amount > e.Amount {
			e = tempExpense
		}
	}
	return e
}

func (es *FileStore) GetExpensesByCategory(filter model.FilterFunc) ([]string, []float64) {
	categories := make(map[string]float64)

	var categoryNames []string
	var categoryValues []float64

	categoryNames = es.categories

	for _, tempExpense := range es.expenses {

		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		categories[tempExpense.Category] += tempExpense.Amount

		if !slices.Contains(categoryNames, tempExpense.Category) {
			categoryNames = append(categoryNames, tempExpense.Category)
		}
	}

	slices.Sort(categoryNames)

	for _, cName := range categoryNames {
		categoryValues = append(categoryValues, categories[cName])
	}

	return categoryNames, categoryValues
}

func (es *FileStore) GetIncomeByCategory(filter model.FilterFunc, userId int64) ([]string, []float64) {
	return []string{"ceva"}, []float64{}
}

func (es *FileStore) GetExpensesByMonth(filter model.FilterFunc) ([]string, []float64) {

	months := map[string]float64{
		"January":   0,
		"February":  0,
		"March":     0,
		"April":     0,
		"May":       0,
		"June":      0,
		"July":      0,
		"August":    0,
		"September": 0,
		"October":   0,
		"November":  0,
		"December":  0,
	}

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}
		if tempExpense.Type == model.INCOME {
			continue
		}

		months[tempExpense.Date.Month().String()] += tempExpense.Amount
	}

	monthIndexes := map[string]float64{
		"January":   1,
		"February":  2,
		"March":     3,
		"April":     4,
		"May":       5,
		"June":      6,
		"July":      7,
		"August":    8,
		"September": 9,
		"October":   10,
		"November":  11,
		"December":  12,
	}

	monthNames := []string{"January", "February", "March", "April",
		"May", "June", "July", "August",
		"September", "October", "November", "December"}

	sort.Slice(monthNames, func(i, j int) bool {
		return monthIndexes[monthNames[i]] < monthIndexes[monthNames[j]]
	})

	var monthValues []float64

	for _, mName := range monthNames {
		monthValues = append(monthValues, months[mName])
	}

	return monthNames, monthValues
}

func (es *FileStore) GetExpensesByDay(filter model.FilterFunc) map[string]float64 {
	days := make(map[string]float64)

	for _, tempExpense := range es.expenses {

		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		day := tempExpense.Date.Day()
		month := tempExpense.Date.Month()
		year := tempExpense.Date.Year()

		key := fmt.Sprintf("%02d/%02d/%04d", day, month, year)

		days[key] += tempExpense.Amount
	}
	return days
}

func (es *FileStore) GetLongestStreakWithoutExpense(filter model.FilterFunc) int {

	var dates []time.Time

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		dates = append(dates, tempExpense.Date)
	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})

	maxStreak := 0

	if len(dates) == 1 {
		return 1
	}

	if len(dates) > 1 {
		firstDate := dates[0]
		for i := 1; i < len(dates); i++ {

			tempStreak := dates[i].Day() - firstDate.Day()

			if tempStreak > maxStreak {
				maxStreak = tempStreak
			}

			firstDate = dates[i]
		}
	}

	return maxStreak
}

func (es *FileStore) GetCountsByCategory(filter model.FilterFunc) ([]string, []float64) {

	categoryFrequencies := make(map[string]float64)

	var categoryNames []string
	var categoryCounts []float64

	categoryNames = es.categories

	for _, tempExpense := range es.expenses {
		if filter(tempExpense) {
			continue
		}

		if tempExpense.Type == model.INCOME {
			continue
		}

		categoryFrequencies[tempExpense.Category] += 1

		if !slices.Contains(categoryNames, tempExpense.Category) {
			categoryNames = append(categoryNames, tempExpense.Category)
		}
	}

	slices.Sort(categoryNames)

	for _, cName := range categoryNames {
		categoryCounts = append(categoryCounts, categoryFrequencies[cName])
	}

	return categoryNames, categoryCounts
}

func (es *FileStore) GetCategories() []string {
	return es.categories
}
