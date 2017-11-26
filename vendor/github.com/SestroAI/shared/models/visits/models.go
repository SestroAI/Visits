package visits

import (
	"github.com/SestroAI/shared/utils"
	"time"
	"github.com/SestroAI/shared/models/orders"
)

type Rating struct {
	ReviewerId string `json:"reviewerId"`
	Value int `json:"value"`
	Comments string `json:"comments"`
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
	Payer 		 string 		   `json:"payer"`
	GuestRating  *Rating 		   `json:"guestRating"`
	MerchantId 	 string 		   `json:"merchantId"`
}

func NewMerchantVisit(id string) *MerchantVisit {
	visit := MerchantVisit{
		IsComplete:   false,
		StartTime:    time.Now(),
		IsOpenForAll: false,
		Payer: "",
		GuestRating:&Rating{},
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
	Orders 		 map[string]*orders.Order   `json:"orders",mapstructure:"orders,squash"`
	List 		 []string 		  `json:"list",mapstructure:",squash"`
	MerchantRating *Rating		  `json:"merchantRating",mapstructure:"merchantRating,squash"`
	Payer		 string 		  `json:"payer"`
}

func NewVisitDinerSession() *VisitDinerSession {
	sess := VisitDinerSession{}
	sess.ID = utils.GenerateUUID()
	sess.Orders = make(map[string]*orders.Order, 0)
	sess.Payer = ""
	sess.MerchantRating = &Rating{}
	sess.DinerId = ""
	return &sess
}