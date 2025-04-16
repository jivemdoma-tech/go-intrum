package gointrum

type StockAttachResponse struct {
	*Response
	Data map[string]*StockAttachData `json:"data,omitempty"`
}

type StockAttachData struct {
	Requests []string `json:"requests"`
}
