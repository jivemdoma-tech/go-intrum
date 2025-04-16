package gointrum

type StockInsertResponse struct {
	*Response
	Data   []uint64 `json:"data,omitempty"`
}
