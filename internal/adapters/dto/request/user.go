package request

type (
	LoginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	GetUserInfo struct {
		ID int
	}

	GenerateToken struct {
		ClientID     string
		ClientSecret string
	}

	CreateUser struct {
		Email    string  `json:"email"`
		Password string  `json:"password"`
		Balance  float64 `json:"balance"`
	}
)
