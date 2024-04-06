package request

type (
	Disburse struct {
		UserID        int     `json:"user_id"`
		Amount        float64 `json:"amount"`
		BankCode      string  `json:"bank_code"`
		AccountNumber string  `json:"account_number"`
		AccountName   string  `json:"account_name"`
	}
	TransferCallback struct {
		TransferID int    `json:"transfer_id"`
		Status     string `json:"status"`
	}
)
