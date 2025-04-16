package gointrum

type SalesUpdateResponse struct {
	*Response
	Data   bool   `json:"data,omitempty"`
}
