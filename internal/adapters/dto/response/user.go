package response

type (
	LoginUser struct {
		AccessToken string `json:"access_token"`
	}

	GenerateToken struct {
		ApiToken string `json:"api_token"`
	}

	CreateUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	GetUserInfo struct {
		ID         int        `json:"id"`
		Email      string     `json:"email"`
		Balance    float64    `json:"balance"`
		Credential Credential `json:"credential"`
	}
)
