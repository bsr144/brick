package entities

import "time"

type User struct {
	ID        int
	Email     string
	Password  string
	Salt      string
	Balance   float64
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (u *User) IsExist() bool {
	return u.ID != 0
}

func (u *User) IsNotExist() bool {
	return u.ID == 0
}
