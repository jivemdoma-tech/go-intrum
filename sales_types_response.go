package intrum

type SalesTypesResponse struct {
	*Response
	Data []*SalesTypesData `json:"data,omitempty"`
}
type SalesTypesData struct {
	ID     int64         `json:"id,string"`
	Name   string        `json:"name"`
	Stages []*SalesStage `json:"stages"`
}
type SalesStage struct {
	ID        int64  `json:"id,string"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Order     int64  `json:"order,string"`
	IsSuccess *bool  `json:"is_success,omitempty"`
	IsFail    *bool  `json:"is_fail,omitempty"`
}
