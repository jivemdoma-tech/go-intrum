package gointrum

type WorkerDepartmentResponse struct {
	*Response
	Data []*WorkerDepartmentData `json:"data"`
}

type WorkerDepartmentData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Publ int64  `json:"publ"`
	// ParentID int64  `json:"parent_id"`       // TODO
	// Order    int64  `json:"order"`           // TODO
	// Description *string `json:"description"` // TODO
	// Timezone string      `json:"timezone"`   // TODO
	// Fields   []interface{} `json:"fields"`   // TODO
}
