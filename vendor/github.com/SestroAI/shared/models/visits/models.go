package visits

import (
	"github.com/SestroAI/shared/models/merchant/menu"
	"github.com/SestroAI/shared/utils"
	"time"
	"github.com/SestroAI/shared/models/orders"
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
	ID           string   		  `json:"id"`
	DinerId      string   		  `json:"dinerId"`
	Orders 		 map[string]*orders.Order   `json:"orders"`
	List 		 []string 		  `json:"list"` 
}

func NewVisitDinerSession() *VisitDinerSession {
	sess := VisitDinerSession{}
	sess.ID = utils.GenerateUUID()
	sess.Orders = make(map[string]*orders.Order, 0)
	return &sess
}