package entities

import "time"

type RecipientAccount struct {
	ID                 int
	AccountNumber      string
	AccountName        string
	BankCode           string
	BankName           string
	VerificationStatus string
	LastVerifiedAt     *time.Time
}

func (r *RecipientAccount) IsExist() bool {
	return r.ID != 0
}

func (r *RecipientAccount) IsNotExist() bool {
	return r.ID == 0
}
