package gointrum

type PurchaserInsertResponse struct {
	*Response
	Data []uint64 `json:"data,omitempty"`
}
