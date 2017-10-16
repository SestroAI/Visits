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
	OrderedAt time.Time `json:"orderedAt",mapstructure:"orderedAt,squash"`
	DeliveredAt time.Time `json:"deliveredAt",mapstructure:"deliveredAt,squash"`
	ItemId string `json:"itemId"`
	Status string `json:"status"`
	Comments string `json:"comments"`
}

func NewOrder() *Order {
	o := Order{
		ID: utils.GenerateUUID(),
		OrderedAt:time.Now(),
		Status:"ordered",
	}
	return &o
}