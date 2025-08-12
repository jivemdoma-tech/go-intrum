package gointrum

type WorkerUpdateResponse struct {
	*Response
	Data bool `json:"data,omitempty"`
}
