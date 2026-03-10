package intrum

type WorkerUpdateResponse struct {
	*Response
	Data bool `json:"data,omitempty"`
}
