package response

type HTTPResponseError struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Error   interface{} `json:"error,omitempty"`
}

func NewHTTPResponseError(code int, errorObj interface{}) *HTTPResponseError {
	return &HTTPResponseError{
		Code:    code,
		Success: false,
		Error:   errorObj,
	}
}
