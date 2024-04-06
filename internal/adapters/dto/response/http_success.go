package response

type HTTPResponseSuccess struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type HTTPResponseSuccessWithPagination struct {
	Code       int         `json:"code"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination"`
}

func NewHTTPResponseSuccess(code int, data interface{}) *HTTPResponseSuccess {
	return &HTTPResponseSuccess{
		Code:    code,
		Success: true,
		Data:    data,
	}
}

func NewHTTPResponseSuccessWithPagination(code int, data interface{}, pagination *PaginationResponse) *HTTPResponseSuccessWithPagination {
	return &HTTPResponseSuccessWithPagination{
		Code:       code,
		Success:    true,
		Data:       data,
		Pagination: pagination,
	}
}
