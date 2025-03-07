package intrumgo

type StockInsertResponse struct {
	Status string  `json:"status"`
	Data   []uint64 `json:"data"`
}
