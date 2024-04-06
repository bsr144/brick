package response

type (
	Credential struct {
		ID           int    `json:"id"`
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}
)
