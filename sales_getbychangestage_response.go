package gointrum

type SalesGetByChangeStageResponse struct {
	Status string                     `json:"status"`
	Data   *SalesGetByChangeStageData `json:"data"`
}

type SalesGetByChangeStageData struct {
	List []*SalesGetByChangeStageDataList `json:"list"`
}

type SalesGetByChangeStageDataList struct {
	SaleID     string `json:"sale_id"`
	SaleTypeID string `json:"sale_type_id"`
	ToStage    string `json:"to_stage"`
	FromStage  string `json:"from_stage"`
	Date       string `json:"date"`
}
