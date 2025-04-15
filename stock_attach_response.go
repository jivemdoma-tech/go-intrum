package gointrum

type StockAttachResponse struct {
	Status string                      `json:"status"`
	Data   map[string]*StockAttachData `json:"data"`
}

type StockAttachData struct {
	Requests []string `json:"requests"`
}
