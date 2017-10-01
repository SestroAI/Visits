package events

import (
	"github.com/SestroAI/shared/firebase/messaging"
	"github.com/SestroAI/shared/models/visits"
)

const (
	VisitEndEventName = "visit_end"
)

type VisitEndEvent struct {
	UId string `json:"uid"`
	Visit *visits.MerchantVisit `json:"visit"`
}

func SendEndVisitEvent(userId string, visit *visits.MerchantVisit) error {
	data := &VisitEndEvent{
		UId:userId,
		Visit:visit,
	}
	return messaging.SendMessage(VisitEndEventName, data)
}
