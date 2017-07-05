package dao

import (
	"github.com/google/logger"
	"errors"

	"github.com/SestroAI/shared/models/transactions"
)

type TransactionDao struct {
	Dao
	BasePath string
}

func NewTransactionDao(token string) *VisitDao {
	return &VisitDao{
		Dao: *NewDao(token),
		BasePath:VISIT_BASE_PATH,
	}
}

func (ref *TransactionDao) SaveTransaction(id string, diner transactions.Transaction) error {
	err := ref.SaveObjectById(id, diner, ref.BasePath)

	if err != nil {
		logger.Errorf("Unable to save Transaction object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *TransactionDao) GetTransaction(id string) (*transactions.Transaction, error) {
	object, _ := ref.GetObjectById(id, ref.BasePath)
	if object == nil {
		return nil, errors.New("Unable to get Transaction with id = " + id)
	}

	transaction := transactions.Transaction{}
	MapToStruct(object.(map[string]interface{}), &transaction)

	return &transaction, nil
}
