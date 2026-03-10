package gointrum

type StockInsertResponse struct {
	Status  string  `json:"status,omitempty"`
	Message string  `json:"message,omitempty"`
	Data    []int64 `json:"data,omitempty"`
}

func (r *StockInsertResponse) GetErrorMessage() string {
	switch {
	case r == nil:
		return ""
	default:
		return ""
	case r.Status != "" && r.Message != "":
		return r.Status + ": " + r.Message
	case r.Message != "":
		return r.Message
	}
}
