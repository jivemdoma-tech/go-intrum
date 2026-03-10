package intrum

type StockDeleteResp struct {
	*Response
	Data bool `json:"data,omitempty"`
}
