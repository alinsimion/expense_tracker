package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/alinsimion/expense_tracker/model"
	_ "modernc.org/sqlite"
)

type SqlLiteStore struct {
	db *sql.DB
}

const (
	sqliteFileName = "database.db"

	expensesTable = "expenses"
	userTable     = "users"
)

var (
	dropUserTable   = fmt.Sprintf(`DROP TABLE %s;`, userTable)
	createUserTable = fmt.Sprintf(`CREATE TABLE %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		avatar_url TEXT
	);`, userTable)

	dropExpenseTable   = fmt.Sprintf(`DROP TABLE %s;`, expensesTable)
	createExpenseTable = fmt.Sprintf(`CREATE TABLE %s(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount int,
		category TEXT,
		currency TEXT,
		date TIMESTAMP,
		description TEXT,
		type TEXT,
		uid INTEGER,
		FOREIGN KEY(uid) REFERENCES users(id) 
	);`, expensesTable)
)

func DropAndCreateTable(db *sql.DB, dropTableQuery string, createTableQuery string) {

	_, err := db.Exec(dropTableQuery)

	if err != nil {
		fmt.Println("No table to drop")
	} else {
		fmt.Println("Table dropped")
	}

	_, err = db.Exec(createTableQuery)

	if err != nil {
		fmt.Println("Could not create table because of error " + err.Error())
	} else {
		fmt.Println("Table created")
	}
}

func NewSqlLiteStore(fileName string) *SqlLiteStore {

	sqliteDatabase, err := sql.Open("sqlite", sqliteFileName) // Open the created SQLite File
	// defer sqliteDatabase.Close()                                     // Defer Closing the database
	if err != nil {
		slog.Error("Error while opening sqlite file", "err", err.Error())
	}
	// CreateExpenseTable(sqliteDatabase)
	// DropAndCreateTable(sqliteDatabase, dropUserTable, createUserTable)
	// DropAndCreateTable(sqliteDatabase, dropExpenseTable, createExpenseTable)

	return &SqlLiteStore{
		db: sqliteDatabase,
	}
}

// --------------------- Users

func (s *SqlLiteStore) GetUserById(id string) *model.User {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 LIMIT 1;", userTable)

	row := s.db.QueryRow(query, id)

	var tempUser model.User

	err := row.Scan(&tempUser.Id, &tempUser.Name, &tempUser.Email, &tempUser.AvatarUrl)
	if err != nil {
		fmt.Println("Error while scaning users", err.Error())
	}

	return &tempUser
}

func (s *SqlLiteStore) GetUserByEmail(email string) *model.User {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1 LIMIT 1;", userTable)

	row := s.db.QueryRow(query, email)

	var tempUser model.User

	err := row.Scan(&tempUser.Id, &tempUser.Name, &tempUser.Email, &tempUser.AvatarUrl)
	if err != nil {
		fmt.Println("Error while scaning users", err.Error())
	}

	return &tempUser
}

func (s *SqlLiteStore) GetUsers(skip int, limit int) []model.User {
	var users []model.User

	var query string

	query = fmt.Sprintf("SELECT * FROM %s ORDER BY id;", userTable)

	rows, err := s.db.Query(query)

	if err != nil {
		fmt.Println("Error while scaning users", err.Error())
	}

	for rows.Next() {
		var tempUser model.User

		err := rows.Scan(&tempUser.Id, &tempUser.Name, &tempUser.Email, &tempUser.AvatarUrl)

		if err != nil {
			fmt.Println("Error while scaning metric", err.Error())
		}

		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows while querying all users")
		case err != nil:
			fmt.Println("cannot retrieve users", err)
		default:

		}
		users = append(users, tempUser)
	}

	return users
}

func (s *SqlLiteStore) CreateUser(e model.User) {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(name, email, avatar_url) 
		VALUES
		($1, $2, $3) RETURNING id`, userTable)

	var lastInsertedId int64
	err := s.db.QueryRow(
		query,
		e.Name,
		e.Email,
		e.AvatarUrl,
	).Scan(&lastInsertedId)

	if err != nil {
		slog.Error("Error while creating User", "err", err.Error())
	}
}

func (s *SqlLiteStore) UpdateUser(id string, e model.User) (int64, error) {
	query := fmt.Sprintf(`
	UPDATE %s
		SET name = $2,
			email = $3,
			avatar_url = $4
    WHERE id = $1`, userTable)

	result, err := s.db.Exec(
		query, id, e.Name, e.Email, e.AvatarUrl,
	)

	if err != nil {
		return -1, fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}

// --------------------- Expenses

func (s *SqlLiteStore) GetExpenseById(id string) *model.Expense {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1 LIMIT 1;", expensesTable)

	row := s.db.QueryRow(query, id)

	var tempExpense model.Expense

	err := row.Scan(&tempExpense.Id, &tempExpense.Amount, &tempExpense.Category, &tempExpense.Currency, &tempExpense.Date, &tempExpense.Description, &tempExpense.Type, &tempExpense.Uid)
	if err != nil {
		fmt.Println("Error while scaning expenses", err.Error())
	}

	return &tempExpense
}
func (s *SqlLiteStore) GetExpenses(skip int, limit int, uid int64) []model.Expense {
	var expenses []model.Expense

	var query string

	query = fmt.Sprintf("SELECT * FROM %s WHERE uid = $1 ORDER BY id;", expensesTable)

	rows, err := s.db.Query(query, uid)

	if err != nil {
		fmt.Println("Error while scaning expenses", err.Error())
	}

	for rows.Next() {
		var tempExpense model.Expense

		err := rows.Scan(&tempExpense.Id, &tempExpense.Amount, &tempExpense.Category, &tempExpense.Currency, &tempExpense.Date, &tempExpense.Description, &tempExpense.Type, &tempExpense.Uid)
		if err != nil {
			fmt.Println("Error while scaning metric", err.Error())
		}

		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No rows while querying all expenses")
		case err != nil:
			fmt.Println("cannot retrieve expenses", err)
		default:

		}
		expenses = append(expenses, tempExpense)
	}

	return expenses
}
func (s *SqlLiteStore) CreateExpense(e model.Expense) {
	query := fmt.Sprintf(`
		INSERT INTO %s
		(amount, category, currency, date, description, type, uid) 
		VALUES
		($1, $2, $3, $4, $5, $6, $7) RETURNING id`, expensesTable)

	var lastInsertedId int64
	err := s.db.QueryRow(
		query,
		e.Amount,
		e.Category,
		e.Currency,
		e.Date,
		e.Description,
		e.Type,
		e.Uid,
	).Scan(&lastInsertedId)

	if err != nil {
		slog.Error("Error while creating expense", "err", err.Error())
		// return -1, fmt.Errorf("failed to execute query: %v", err)
	}
	// return lastInsertedId, nil
}
func (s *SqlLiteStore) UpdateExpense(id string, e model.Expense) (int64, error) {
	query := fmt.Sprintf(`
	UPDATE %s
		SET amount = $2,
			category = $3,
			currency = $4,
			date = $5,
			description = $6,
			type = $7
			uid = $8
    WHERE id = $1`, expensesTable)

	result, err := s.db.Exec(
		query, id, e.Amount, e.Category, e.Currency, e.Date, e.Description, e.Type, e.Uid,
	)

	if err != nil {
		return -1, fmt.Errorf("failed to execute query: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()

	return rowsAffected, nil
}
func (s *SqlLiteStore) DeleteExpense(id string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
		RETURNING id`, expensesTable)

	_, err := s.db.Exec(
		query, id,
	)

	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return nil
}

func (s *SqlLiteStore) GetCurrentBalance(filter model.FilterFunc, userId int64) float64 {
	expenses := s.GetExpenses(0, 0, userId)

	var balance float64
	for _, expense := range expenses {
		if filter(expense) {
			continue
		}
		if expense.Type == model.INCOME {
			balance += expense.Amount
		} else {
			balance -= expense.Amount
		}
	}
	return balance
}

func (s *SqlLiteStore) GetLargestExpense(filter model.FilterFunc, userId int64) model.Expense {
	expenses := s.GetExpenses(0, 0, userId)

	var maxExpense model.Expense
	for _, expense := range expenses {
		if filter(expense) {
			continue
		}
		if expense.Amount > maxExpense.Amount {
			maxExpense = expense
		}
	}
	return maxExpense
}

func (s *SqlLiteStore) GetExpensesByCategory(filter model.FilterFunc, userId int64) ([]string, []float64) {
	categories := make(map[string]float64)

	expenses := s.GetExpenses(0, 0, userId)

	for _, expense := range expenses {
		if filter(expense) {
			continue
		}
		categories[expense.Category] += expense.Amount
	}

	var categoryNames []string
	var categoryTotals []float64

	for key, value := range categories {
		categoryNames = append(categoryNames, key)
		categoryTotals = append(categoryTotals, value)
	}

	return categoryNames, categoryTotals
}

func (s *SqlLiteStore) GetExpensesByMonth(filter model.FilterFunc, userId int64) ([]string, []float64) {
	months := make(map[string]float64)

	expenses := s.GetExpenses(0, 0, userId)

	for _, expense := range expenses {
		if expense.Type != model.EXPENSE {
			continue
		}
		if filter(expense) {
			continue
		}
		month := expense.Date.Month().String()
		months[month] += expense.Amount
	}

	var monthNames []string
	var monthTotals []float64

	for key, value := range months {
		monthNames = append(monthNames, key)
		monthTotals = append(monthTotals, value)
	}

	return monthNames, monthTotals
}
func (s *SqlLiteStore) GetExpensesByDay(filter model.FilterFunc, userId int64) map[string]float64 {
	days := make(map[string]float64)

	expenses := s.GetExpenses(0, 0, userId)

	for _, expense := range expenses {

		if expense.Type != model.EXPENSE {
			continue
		}
		if filter(expense) {
			continue
		}
		day := expense.Date.Day()
		month := expense.Date.Month()
		year := expense.Date.Year()

		key := fmt.Sprintf("%02d/%02d/%04d", day, month, year)

		days[key] += expense.Amount
	}

	var categoryNames []string
	var categoryCounts []float64

	for key, value := range days {
		categoryNames = append(categoryNames, key)
		categoryCounts = append(categoryCounts, value)
	}

	return days
}

func (s *SqlLiteStore) GetLongestStreakWithoutExpense(filter model.FilterFunc, userId int64) int {
	return 5
}
func (s *SqlLiteStore) GetCountsByCategory(filter model.FilterFunc, userId int64) ([]string, []float64) {
	categories := make(map[string]float64)

	expenses := s.GetExpenses(0, 0, userId)

	for _, expense := range expenses {
		if filter(expense) {
			continue
		}
		categories[expense.Category] += 1.0
	}

	var categoryNames []string
	var categoryCounts []float64

	for key, value := range categories {
		categoryNames = append(categoryNames, key)
		categoryCounts = append(categoryCounts, value)
	}

	return categoryNames, categoryCounts

}
func (s *SqlLiteStore) GetCurrentIncomes(filter model.FilterFunc, userId int64) float64 {
	expenses := s.GetExpenses(0, 0, userId)
	var total float64
	for _, expense := range expenses {
		if filter(expense) {
			continue
		}
		if expense.Type == model.INCOME {
			total += expense.Amount
		}
	}
	return total
}
func (s *SqlLiteStore) GetCurrentExpenses(filter model.FilterFunc, userId int64) float64 {
	expenses := s.GetExpenses(0, 0, userId)

	var total float64
	for _, expense := range expenses {
		if filter(expense) {
			continue
		}
		if expense.Type == model.EXPENSE {
			total += expense.Amount
		}
	}
	return total
}

func (s *SqlLiteStore) GetCategories(userId int64) []string {
	expenses := s.GetExpenses(0, 0, userId)

	var categories []string
	for _, expense := range expenses {

		if !StringSliceContains(categories, expense.Category) {
			categories = append(categories, expense.Category)
		}
	}
	return categories
}

func StringSliceContains(slice []string, needle string) bool {
	for _, elem := range slice {
		if elem == needle {
			return true
		}
	}
	return false
}
