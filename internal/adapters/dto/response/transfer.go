package response

type (
	ValidateAccount struct {
		BankCode      string `json:"bank_code"`
		BankName      string `json:"bank_name"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
	}

	ValidateBank struct {
		Success       bool   `json:"success"`
		BankCode      string `json:"bank_code"`
		BankName      string `json:"bank_name"`
		AccountNumber string `json:"account_number"`
		AccountName   string `json:"account_name"`
	}
	Disburse struct {
		TransferID    int     `json:"transfer_id"`
		Status        string  `json:"status"`
		Amount        float64 `json:"amount"`
		BankCode      string  `json:"bank_code"`
		BankName      string  `json:"bank_name"`
		AccountNumber string  `json:"account_number"`
		AccountName   string  `json:"account_name"`
		CreatedAt     string  `json:"created_at"`
	}
)
