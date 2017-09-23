package dao

import (
	"github.com/SestroAI/shared/models/merchant"
	"github.com/SestroAI/shared/models/visits"
	"errors"
	"fmt"
)

type RestaurantDao struct {
	Dao
}

const (
	RESTAURANT_PATH = "/restaurants"
	TABLE_PATH = "/tables"
)


func NewRestaurantDao(token string) *RestaurantDao {
	return &RestaurantDao{
		Dao: *NewDao(token),
	}
}

func (ref *RestaurantDao) SaveRestaurant(id string, restro *merchant.Merchant) error {
	err := ref.SaveObjectById(id, restro, RESTAURANT_PATH)

	if err != nil {
		return err
	}

	return nil
}

func (ref *RestaurantDao) GetRestaurantById(id string) (*merchant.Merchant, error) {
	object, err := ref.GetObjectById(id, RESTAURANT_PATH)
	if err != nil || object == nil{
		return nil, err
	}

	restro := merchant.Merchant{}

	err = MapToStruct(object.(map[string]interface{}), &restro)

	fmt.Println("restro = ", restro, object)

	return &restro, err
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
	err := ref.SaveObjectById(id, table, TABLE_PATH)

	if err != nil {
		return err
	}

	return nil
}

func (ref *RestaurantDao) UpdateTableOngoingVisit(tableId string, visit *visits.MerchantVisit) error {
	table, err := ref.GetTableById(tableId)
	if err != nil{
		return err
	}

	table.OngoingVisitId = visit.ID
	err = ref.SaveTable(table.ID, table)
	if err != nil {
		return err
	}
	return nil
}
