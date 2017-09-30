package orders

import (
	"time"
	"github.com/SestroAI/shared/utils"
)

var AllowedOrderStatus = []string{
	"ordered",
	"cooking",
	"delivered",
}

type Order struct {
	ID string `json:"id"`
	SessionID string `json:"sessionId"`
	OrderedAt time.Time `json:"orderedAt"`
	ServedAt time.Time `json:"deliveredAt"`
	ItemId string `json:"itemId"`
	Status string `json:"status"`
}

func NewOrder() *Order {
	o := Order{
		ID: utils.GenerateUUID(),
		OrderedAt:time.Now(),
		Status:"ordered",
	}
	return &o
}