package gointrum

type WorkerFiliationResponse struct {
	*Response
	Data []*WorkerFiliationData `json:"data"`
}

type WorkerFiliationData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Description interface{} `json:"description"` // TODO
	// City        string      `json:"city"`        // TODO
	// Adress      string      `json:"adress"`      // TODO
	// Phone       string      `json:"phone"`       // TODO
	// Email       string      `json:"email"`       // TODO
}
