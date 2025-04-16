package gointrum

// Ответ API Intrum
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func (r *Response) GetErrorMessage() string {
	return r.Message
}
