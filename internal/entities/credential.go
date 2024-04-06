package entities

import "time"

type Credential struct {
	ID           int
	ClientID     string
	ClientSecret string
	UserID       int
	CreatedAt    *time.Time
}

func (c *Credential) IsNotExist() bool {
	return c.ID == 0
}
