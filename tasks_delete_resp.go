package gointrum

type TasksDeleteResp struct {
	*Response
	Data []any `json:"data"`
}
