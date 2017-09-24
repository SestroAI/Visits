package transactions

import "time"

type Transaction struct {
	ID                 string `json:"id"`
	Time               time.Time
	VisitId            string `json:"visitId"`
	ConfirmationNumber int    `json:"confirmationNumber"`
}
