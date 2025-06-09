package gointrum

type SalesTypesResponse struct {
	*Response
	Data []*SalesTypesData `json:"data,omitempty"`
}
type SalesTypesData struct {
	ID     uint32        `json:"id,string"`
	Name   string        `json:"name"`
	Stages []*SalesStage `json:"stages"`
}
type SalesStage struct {
	ID        uint32 `json:"id,string"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Order     uint32 `json:"order,string"`
	IsSuccess *bool  `json:"is_success,omitempty"`
	IsFail    *bool  `json:"is_fail,omitempty"`
}
