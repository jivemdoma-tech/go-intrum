package gointrum

type PurchaserInsertResp struct {
	*Response
	Data []int64 `json:"data,omitempty"`
}
