package gointrum

type PurchaserInsertResp struct {
	*Response
	Data []uint64 `json:"data,omitempty"`
}
