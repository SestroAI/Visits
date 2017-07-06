package dao

import (
	"github.com/google/logger"
	"errors"

	"github.com/SestroAI/shared/models/visits"
	"github.com/SestroAI/shared/models/auth"
)

const (
	VISIT_BASE_PATH = "/visits"
	VISIT_SESSION_RELATIVE_PATH = "/sessions"
)

type VisitDao struct {
	Dao
	BasePath string
}

func NewVisitDao(token string) *VisitDao {
	return &VisitDao{
		Dao: *NewDao(token),
		BasePath:VISIT_BASE_PATH,
	}
}

func (ref *VisitDao) StartNewVisit(diner *auth.Diner, tableId string) (*visits.RestaurantVisit, error) {

	//Check if table is empty
	restaurantDao := NewRestaurantDao(ref.Token)
	table, err := restaurantDao.GetTableById(tableId)
	if err != nil {
		logger.Errorf("No table with tableId = %s exists to start visit for Diner with id = %s", tableId, diner.ID)
		return nil, err
	}
	if table.OngoingVisitId != "" {
		logger.Errorf("Cannot start a new visit on tableId = %s for diner id = %s when there is already an active visit going on", tableId, diner.ID)
		return nil, err
	}

	visit := visits.NewRestaurantVisit("")
	visit.TableId = tableId

	visitorSession := visits.NewVisitDinerSession("")
	visitorSession.DinerId = diner.ID
	err = ref.SaveVisitSession(visitorSession.ID, visitorSession); if err != nil {
		return nil, err
	}

	visit.Diners[diner.ID] = visitorSession
	err = ref.SaveVisit(visit.ID, visit); if err != nil {
		logger.Errorf("Unable to save Visit object (create new visit) with tableId = %s", tableId)
		return nil, err
	}

	err = restaurantDao.UpdateTableOngoingVisit(table.ID, visit)
	if err != nil {
		logger.Errorf("Unable to update table ongoing visit for tableId = %s and visit id = %s", tableId, visit.ID)
		return nil, err
	}

	userDao := NewUserDao(ref.Token)
	err = userDao.UpdateDinerOngoingVisit(diner.ID, visit)
	if err != nil{
		logger.Errorf("Unable to update visit id = %s for diner id = %s", visit.ID, diner.ID)
		return nil, err
	}

	return visit, nil
}


func (ref *VisitDao) SaveVisit(id string, visit *visits.RestaurantVisit) error {
	err := ref.SaveObjectById(id, visit, ref.BasePath)

	if err != nil {
		logger.Errorf("Unable to save Visit object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *VisitDao) GetVisit(id string) (*visits.RestaurantVisit, error) {
	object, _ := ref.GetObjectById(id, ref.BasePath)
	if object == nil {
		return nil, errors.New("Unable to get Visit with id = " + id)
	}

	visit := visits.RestaurantVisit{}
	MapToStruct(object.(map[string]interface{}), &visit)

	return &visit, nil
}

func (ref *VisitDao) SaveVisitSession(id string, visitSess *visits.VisitDinerSession) error {
	err := ref.SaveObjectById(id, *visitSess, ref.BasePath + VISIT_SESSION_RELATIVE_PATH)

	if err != nil {
		logger.Errorf("Unable to save Visit Session object with Id = %s", id)
		return err
	}

	return nil
}

func (ref *VisitDao) GetVisitSession(id string) (*visits.VisitDinerSession, error) {
	object, _ := ref.GetObjectById(id, ref.BasePath)
	if object == nil {
		return nil, errors.New("Unable to get Visit Session with id = " + id)
	}

	visitSess := visits.VisitDinerSession{}
	MapToStruct(object.(map[string]interface{}), &visitSess)

	return &visitSess, nil
}