package intrum

type PurchaserUpdateResponse struct {
	*Response
	Data bool `json:"data,omitempty"`
}
