package intrum

type PurchaserInsertResp struct {
	*Response
	Data []int64 `json:"data,omitempty"`
}
