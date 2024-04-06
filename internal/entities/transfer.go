package entities

import "time"

type Transfer struct {
	ID                 int
	RecipientAccountID int
	SenderAccountID    int
	Amount             float64
	Status             string
	CreatedAt          *time.Time
	CompletedAt        *time.Time
}
