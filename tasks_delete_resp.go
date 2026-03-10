package intrum

type TasksDeleteResp struct {
	*Response
	Data []any `json:"data"`
}
