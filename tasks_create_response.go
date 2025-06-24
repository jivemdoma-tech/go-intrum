package gointrum

type TasksCreateResponse struct {
	*Response
	Data TasksData `json:"data"`
}

type TasksData struct {
	ID int64 `json:"id"`
}
