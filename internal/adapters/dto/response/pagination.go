package response

type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalCount   int `json:"total_count"`
	TotalPages   int `json:"total_pages"`
	FirstPage    int `json:"first_page"`
	NextPage     int `json:"next_page"`
	LastPage     int `json:"last_page"`
	CurrentCount int `json:"current_count"`
}
