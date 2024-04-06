package request

type (
	Query struct {
		Search string `form:"search"`
		Page   int    `form:"page"`
		Size   int    `form:"size"`
	}

	ValidateAccount struct {
		BankCode      string `form:"bankCode"`
		AccountNumber string `form:"accountNumber"`
	}
)
