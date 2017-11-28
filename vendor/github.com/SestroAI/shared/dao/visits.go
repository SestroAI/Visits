package dao

import (
	"errors"
	"github.com/SestroAI/shared/logger"

	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/models/visits"
	serrors "github.com/SestroAI/shared/utils/errors"
	"time"
)

const (
	VISIT_PATH         = "/visits"
	VISIT_SESSION_PATH = "/visit_sessions"
)

type VisitDao struct {
	Dao
}

func NewVisitDao(token string) *VisitDao {
	return &VisitDao{
		Dao: *NewDao(token),
	}
}

func (ref *VisitDao) StartNewVisit(diner *auth.User, tableId string) (*visits.MerchantVisit, error) {
	//Check if table is empty
	restaurantDao := NewRestaurantDao(ref.Token)
	table, err := restaurantDao.GetTableById(tableId)
	if err != nil {
		logger.Errorf("No table with tableId = %s exists to start visit for Diner with id = %s and "+
			"err = %s", tableId, diner.ID, err.Error())
		return nil, serrors.ErrConflict
	}
	if table.OngoingVisitId != "" {
		logger.Infof("Cannot start a new visit on tableId = %s for diner id = %s when there is already an "+
			"active visit going on", tableId, diner.ID)
		visit, err := ref.GetVisit(table.OngoingVisitId)
		if err != nil {
			return nil, err
		}
		return visit, errors.New("Active Visit for this table already exists")
	}

	visit := visits.NewMerchantVisit("")
	visit.TableId = tableId
	visit.MerchantId = table.MerchantId

	visitorSession := visits.NewVisitDinerSession()
	visitorSession.DinerId = diner.ID
	visitorSession.Payer = diner.ID
	visitorSession.VisitId = visit.ID
	err = ref.SaveVisitSession(visitorSession.ID, visitorSession)
	if err != nil {
		return nil, err
	}

	visit.Diners[diner.ID] = visitorSession.ID
	err = ref.SaveVisit(visit.ID, visit)
	if err != nil {
		return nil, err
	}

	err = restaurantDao.UpdateTableOngoingVisit(table.ID, visit)
	if err != nil {
		return nil, err
	}

	userDao := NewUserDao(ref.Token)
	err = userDao.UpdateDinerOngoingVisit(diner.ID, visit)
	if err != nil {
		return nil, err
	}

	return visit, nil
}

func (ref *VisitDao) EndVisit(visit *visits.MerchantVisit) (error) {
	visit.IsComplete = true
	visit.EndTime = time.Now()

	err := ref.SaveVisit(visit.ID, visit)
	if err != nil {
		return err
	}

	restaurantDao := NewRestaurantDao(ref.Token)
	err = restaurantDao.UpdateTableOngoingVisit(visit.TableId, visit)
	if err != nil {
		return err
	}

	//Need to use admin token to update diner profiles
	userDao := NewUserDao("")
	userDao.IsService = true

	for dinerId, _ := range visit.Diners {
		err = userDao.UpdateDinerOngoingVisit(dinerId, nil)
		if err != nil {
			logger.Errorf("Unable to update ongoing visit for user ID = %s and err = %s", dinerId, err.Error())
			/*
			TODO: Handle this error gracefully
			 */
		}
	}
	return nil
}

func (ref *VisitDao) SaveVisit(id string, visit *visits.MerchantVisit) error {
	return ref.SaveObjectById(id, visit, VISIT_PATH)
}

func (ref *VisitDao) GetVisit(id string) (*visits.MerchantVisit, error) {
	object, err := ref.GetObjectById(id, VISIT_PATH)
	if object == nil || err != nil {
		return nil, errors.New("Unable to get Visit with id = " + id)
	}

	visit := visits.MerchantVisit{}
	MapToStruct(object.(map[string]interface{}), &visit)

	return &visit, nil
}

func (ref *VisitDao) SaveVisitSession(id string, visitSess *visits.VisitDinerSession) error {
	return ref.SaveObjectById(id, *visitSess, VISIT_SESSION_PATH)
}

func (ref *VisitDao) GetVisitSession(id string) (*visits.VisitDinerSession, error) {
	object, err := ref.GetObjectById(id, VISIT_SESSION_PATH)
	if object == nil{
		return nil, errors.New("Unable to get Visit Session with id = " + id)
	}
	if err != nil {
		return nil, err
	}

	visitSess := visits.VisitDinerSession{}
	err = MapToStruct(object.(map[string]interface{}), &visitSess)
	if err != nil {
		return nil, err
	}

	return &visitSess, nil
}
