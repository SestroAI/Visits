package dao

import (
	"github.com/SestroAI/shared/models/restaurant"
	"github.com/SestroAI/shared/models/visits"
	"errors"
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

func (ref *RestaurantDao) SaveRestaurant(id string, restro *restaurant.Restaurant) error {
	err := ref.SaveObjectById(id, restro, RESTAURANT_PATH)

	if err != nil {
		return err
	}

	return nil
}

func (ref *RestaurantDao) GetRestaurantById(id string) (*restaurant.Restaurant, error) {
	object, _ := ref.GetObjectById(id, RESTAURANT_PATH)
	if object == nil {
		return nil, errors.New("Unable to get diner with id = " + id)
	}

	restro := restaurant.Restaurant{}
	MapToStruct(object.(map[string]interface{}), &restro)
	return &restro, nil
}

func (ref *RestaurantDao) GetTableById(id string) (*restaurant.Table, error) {
	object, _ := ref.GetObjectById(id, TABLE_PATH)
	if object == nil {
		return nil, errors.New("Unable to get diner with id = " + id)
	}

	table := restaurant.Table{}
	MapToStruct(object.(map[string]interface{}), &table)
	return &table, nil
}

func (ref *RestaurantDao) SaveTable(id string, table *restaurant.Table) error {
	err := ref.SaveObjectById(id, table, TABLE_PATH)

	if err != nil {
		return err
	}

	return nil
}

func (ref *RestaurantDao) UpdateTableOngoingVisit(tableId string, visit *visits.RestaurantVisit) error {
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
