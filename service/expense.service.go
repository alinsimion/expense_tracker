package service

import (
	"github.com/alinsimion/expense_tracker/db"
	"github.com/alinsimion/expense_tracker/model"
)

func NewExpenseService(db db.DB) *ExpenseService {
	return &ExpenseService{

		ExpenseDB: db,
	}
}

type ExpenseService struct {
	ExpenseDB db.DB
}

func (es *ExpenseService) GetExpense(id string) model.Expense {
	return es.ExpenseDB.Expenses[0]
}

func (es *ExpenseService) GetExpenses(skip int, limit int) []model.Expense {
	return es.ExpenseDB.Expenses[skip : skip+limit]
}

func (es *ExpenseService) AddExpense(e model.Expense) {
	es.ExpenseDB.Expenses = append(es.ExpenseDB.Expenses, e)
}

func (es *ExpenseService) UpdateExpense(id string, e model.Expense) *model.Expense {
	for i := 0; i < len(es.ExpenseDB.Expenses); i++ {
		if es.ExpenseDB.Expenses[i].Description == id {
			es.ExpenseDB.Expenses[i] = e
			return &e
		}
	}
	return nil
}
