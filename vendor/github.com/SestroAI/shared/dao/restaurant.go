package dao

import (
	"github.com/SestroAI/shared/models/restaurant"
	"github.com/google/logger"
	"github.com/SestroAI/shared/models/visits"
)

type RestaurantDao struct {
	Dao
	BasePath string
}

const (
	RESTAURANT_BASE_PATH = "/restaurants"
	RESTAURANT_TABLE_RELATIVE_PATH = "/tables"
)


func NewRestaurantDao(token string) *RestaurantDao {
	return &RestaurantDao{
		Dao: *NewDao(token),
		BasePath:RESTAURANT_BASE_PATH,
	}
}

func (ref *RestaurantDao) SaveRestaurant(id string, restro *restaurant.Restaurant) error {
	err := ref.SaveObjectById(id, restro, ref.BasePath)

	if err != nil {
		logger.Errorf("Unable to save RESTAURANT object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *RestaurantDao) GetRestaurantById(id string) (*restaurant.Restaurant, error) {
	object, err := ref.GetObjectById(id, ref.BasePath)
	if object == nil {
		return nil, err
	}

	restro := restaurant.Restaurant{}
	MapToStruct(object.(map[string]interface{}), &restro)
	return &restro, nil
}

func (ref *RestaurantDao) GetTableById(id string) (*restaurant.Table, error) {
	object, err := ref.GetObjectById(id, ref.BasePath + RESTAURANT_TABLE_RELATIVE_PATH)
	if object == nil || err != nil {
		return nil, err
	}

	table := restaurant.Table{}
	MapToStruct(object.(map[string]interface{}), &table)
	return &table, nil
}

func (ref *RestaurantDao) SaveTable(id string, table *restaurant.Table) error {
	err := ref.SaveObjectById(id, table, ref.BasePath + RESTAURANT_TABLE_RELATIVE_PATH)

	if err != nil {
		logger.Errorf("Unable to save Table object with Id = %s", id)
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
