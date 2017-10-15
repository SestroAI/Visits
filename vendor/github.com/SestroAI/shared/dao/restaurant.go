package dao

import (
	"errors"
	"github.com/SestroAI/shared/models/merchant"
	"github.com/SestroAI/shared/models/visits"
	"github.com/SestroAI/shared/models/merchant/menu"
)

type RestaurantDao struct {
	Dao
}

const (
	RESTAURANT_PATH = "/restaurants"
	TABLE_PATH      = "/tables"
	MERCHANT_STRIPE_INFO_PATH = "/merchant_stripe_info"
)

func NewRestaurantDao(token string) *RestaurantDao {
	return &RestaurantDao{
		Dao: *NewDao(token),
	}
}

func (ref *RestaurantDao) SaveRestaurant(id string, restro *merchant.Merchant) error {
	return ref.SaveObjectById(id, restro, RESTAURANT_PATH)
}

func (ref *RestaurantDao) GetRestaurantById(id string) (*merchant.Merchant, error) {
	object, err := ref.GetObjectById(id, RESTAURANT_PATH)
	if err != nil || object == nil {
		return nil, err
	}

	restro := merchant.Merchant{}

	err = MapToStruct(object.(map[string]interface{}), &restro)

	return &restro, err
}

func (ref *RestaurantDao) GetMerchantStripeInfo(merchantId string) (*merchant.MerchantStripeInfo, error) {
	object, err := ref.GetObjectById(merchantId, MERCHANT_STRIPE_INFO_PATH)
	if err != nil || object == nil {
		return nil, errors.New("Unable to get merchant stripe info with id = " + merchantId)
	}

	info := merchant.MerchantStripeInfo{}
	err = MapToStruct(object, &info)

	return &info, err
}

func (ref *RestaurantDao) SaveMerchantStripeInfo(id string, info *merchant.MerchantStripeInfo) error {
	return ref.SaveObjectById(id, info, MERCHANT_STRIPE_INFO_PATH)
}

func (ref *RestaurantDao) GetTableById(id string) (*merchant.Table, error) {
	object, err := ref.GetObjectById(id, TABLE_PATH)
	if err != nil || object == nil {
		return nil, errors.New("Unable to get restaurant with id = " + id)
	}

	table := merchant.Table{}
	err = MapToStruct(object.(map[string]interface{}), &table)

	return &table, err
}

func (ref *RestaurantDao) SaveTable(id string, table *merchant.Table) error {
	return ref.SaveObjectById(id, table, TABLE_PATH)
}

func (ref *RestaurantDao) UpdateTableOngoingVisit(tableId string, visit *visits.MerchantVisit) error {
	table, err := ref.GetTableById(tableId)
	if err != nil {
		return err
	}

	table.OngoingVisitId = visit.ID
	err = ref.SaveTable(table.ID, table)
	if err != nil {
		return err
	}
	return nil
}

func (ref *RestaurantDao) GetMenuItemById(id string) (*menu.Item, error) {
	object, err := ref.GetObjectById(id, TABLE_PATH)
	if err != nil || object == nil {
		return nil, errors.New("Unable to get item with id = " + id)
	}

	item := menu.Item{}
	err = MapToStruct(object.(map[string]interface{}), &item)

	return &item, err
}