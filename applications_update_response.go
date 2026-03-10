package intrum

type ApplicationsUpdateResponse struct {
	*Response
	Data bool `json:"data,omitempty"`
}
