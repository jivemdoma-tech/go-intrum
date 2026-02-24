package gointrum

type StockInsertResponse struct {
	*Response
	Data []int64 `json:"data,omitempty"`
}
