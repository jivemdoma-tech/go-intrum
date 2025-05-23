package gointrum

type ApplicationsUpdateResponse struct {
	*Response
	Data bool `json:"data,omitempty"`
}
