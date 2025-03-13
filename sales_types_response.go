package gointrum

type SalesTypesResponse struct {
	Status string            `json:"status"`
	Data   []*SalesTypesData `json:"data"`
}
type SalesTypesData struct {
	ID     uint16        `json:"id,string"`
	Name   string        `json:"name"`
	Stages []*SalesStage `json:"stages"`
}
type SalesStage struct {
	ID        uint16 `json:"id,string"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Order     uint16 `json:"order,string"`
	IsSuccess *bool   `json:"is_success,omitempty"`
	IsFail    *bool   `json:"is_fail,omitempty"`
}
