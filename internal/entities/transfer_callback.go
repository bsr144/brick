package entities

import "time"

type TransferCallback struct {
	ID         int
	TransferID int
	Status     string
	ReceivedAt *time.Time
}
