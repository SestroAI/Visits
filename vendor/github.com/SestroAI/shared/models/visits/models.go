package visits

import (
	"github.com/SestroAI/shared/models/merchant/menu"
	"github.com/SestroAI/shared/utils"
	"time"
)

type VisitRating struct {
	menu.Rating
	VisitId string
}

type MerchantVisit struct {
	ID           string            `json:"id"`
	Diners       map[string]string `json:"diners"`
	StartTime    time.Time         `json:"startTime"`
	EndTime      time.Time         `json:"endTime"`
	TableId      string            `json:"tableId"`
	Transactions []string          `json:"transactions"`
	IsComplete   bool              `json:"isComplete"`
	IsOpenForAll bool              `json:"isOpenForAll"`
}

func NewMerchantVisit(id string) *MerchantVisit {
	visit := MerchantVisit{
		IsComplete:   false,
		StartTime:    time.Now(),
		IsOpenForAll: false,
	}
	if id == "" {
		id = utils.GenerateUUID()
	}
	visit.ID = id
	visit.Diners = map[string]string{}
	visit.Transactions = make([]string, 0)
	visit.StartTime = time.Now()
	return &visit
}

type VisitDinerSession struct {
	ID           string   `json:"id"`
	DinerId      string   `json:"dinerId"`
	ItemsInCart  []string `json:"itemsInCart"`
	ItemsOrdered []string `json:"itemsOrdered"`
	ItemsServed  []string `json:"itemsServed"`
}

func NewVisitDinerSession() *VisitDinerSession {
	sess := VisitDinerSession{}
	sess.ID = utils.GenerateUUID()
	sess.ItemsOrdered = make([]string, 0)
	sess.ItemsInCart = make([]string, 0)
	sess.ItemsServed = make([]string, 0)
	return &sess
}
