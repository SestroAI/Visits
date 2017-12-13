package dao

import (
	"github.com/SestroAI/shared/models/billing"
)

type BillingDao struct {
	Dao
}

const (
	USER_BILLS_PATH = "/user_bills"
)

func NewBillingDao(token string) *BillingDao {
	return &BillingDao{
		Dao: *NewDao(token),
	}
}

func (ref *BillingDao) SaveUserBill(id string, bill *billing.UserBill) error {
	return ref.SaveObjectById(id, bill, USER_BILLS_PATH + "/" + bill.BilledToUserId)
}

func (ref *BillingDao) GetUserBillById(id string, userId string) (*billing.UserBill, error) {
	object, err := ref.GetObjectById(id, USER_BILLS_PATH + "/" + userId)
	if err != nil || object == nil {
		return nil, err
	}

	ub := billing.UserBill{}

	err = MapToStruct(object.(map[string]interface{}), &ub)

	return &ub, err
}

func (ref *BillingDao) GetAllUserBills(userId string) ([]*billing.UserBill, error) {
	object, err := ref.GetObjectById(userId, USER_BILLS_PATH)
	if err != nil {
		return nil, err
	}

	list := make([]*billing.UserBill, 0)

	if object == nil {
		//return empty list
		return list, nil
	}

	var billList map[string]billing.UserBill

	err = MapToStruct(object.(map[string]interface{}), &billList)
	if err != nil {
		return nil, err
	}

	//Convert to map to list
	for _, bill := range billList {
		list = append(list, &bill)
	}

	return list, err
}