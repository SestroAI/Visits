package dao

import (
	"errors"
	"github.com/SestroAI/shared/logger"

	"github.com/SestroAI/shared/models/transactions"
)

const (
	TRANSACTION_PATH = "/transactions"
)

type TransactionDao struct {
	Dao
}

func NewTransactionDao(token string) *TransactionDao {
	return &TransactionDao{
		Dao: *NewDao(token),
	}
}

func (ref *TransactionDao) SaveTransaction(id string, transaction transactions.Transaction) error {
	err := ref.SaveObjectById(id, transaction, TRANSACTION_PATH)

	if err != nil {
		logger.Errorf("Unable to save Transaction object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *TransactionDao) GetTransaction(id string) (*transactions.Transaction, error) {
	object, _ := ref.GetObjectById(id, TRANSACTION_PATH)
	if object == nil {
		return nil, errors.New("Unable to get Transaction with id = " + id)
	}

	transaction := transactions.Transaction{}
	MapToStruct(object.(map[string]interface{}), &transaction)

	return &transaction, nil
}
