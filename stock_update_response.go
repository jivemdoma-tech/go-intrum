package intrum

type StockUpdateResponse struct {
	*Response
	Data bool `json:"data,omitempty"`
}
