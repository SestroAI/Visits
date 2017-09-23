package visits

import (
	"time"
	"github.com/SestroAI/shared/utils"
	"github.com/SestroAI/shared/models/merchant/menu"
)

type VisitRating struct {
	menu.Rating
	VisitId string
}

type MerchantVisit struct {
	ID string
	Diners map[string]*VisitDinerSession
	StartTime time.Time `json:"startTime"`
	EndTime time.Time `json:"endTime"`
	TableId string `json:"table"`
	Transactions []string
	IsComplete bool
}

func NewMerchantVisit(id string) *MerchantVisit {
	visit := MerchantVisit{IsComplete:false, StartTime:time.Now()}
	if id == "" {
		id = utils.GenerateUUID()
	}
	visit.ID = id
	visit.Diners = map[string]*VisitDinerSession{}
	visit.Transactions = make([]string, 0)
	return &visit
}

type VisitDinerSession struct {
	ID string `json:"id"`
	DinerId string `json:"dinerId"`
	ItemsInCart []string `json:"itemsInCart"`
	ItemsOrdered []string `json:"itemsOrdered"`
	ItemsServed []string `json:"itemsServed"`
}

func NewVisitDinerSession(id string) *VisitDinerSession {
	sess := VisitDinerSession{}
	if id == "" {
		id = utils.GenerateUUID()
	}
	sess.ID = id
	sess.ItemsOrdered = make([]string, 0)
	sess.ItemsInCart = make([]string, 0)
	sess.ItemsServed = make([]string, 0)
	return &sess
}