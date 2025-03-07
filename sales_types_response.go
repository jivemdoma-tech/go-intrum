package intrumgo

type SalesTypesResponse struct {
	Status string            `json:"status"`
	Data   []*SalesTypesData `json:"data"`
}
type SalesTypesData struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Stages []*SalesStage `json:"stages"`
}
type SalesStage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Order     string `json:"order"`
	IsSuccess *bool  `json:"is_success,omitempty"`
	IsFail    *bool  `json:"is_fail,omitempty"`
}
