package gointrum

import (
	"context"
	"fmt"
	"time"
)

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter-stage-period
type SalesGetByChangeStageParams struct {
	DateStart time.Time
	DateEnd   time.Time
	SaleID    []uint64
	Stage     []uint64
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter-stage-period
func SalesGetByChangeStage(ctx context.Context, subdomain, apiKey string, inputParams *SalesGetByChangeStageParams) (*SalesGetByChangeStageResponse, error) {
	methodURL := fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/getbychangestage", subdomain)

	// Параметры запроса

	params := make(map[string]string, 2+len(inputParams.SaleID)+len(inputParams.Stage))

	// date_start
	if !inputParams.DateStart.IsZero() {
		params["params[date_start]"] = inputParams.DateStart.Format(dateLayout)
	}
	// date_end
	if !inputParams.DateEnd.IsZero() {
		params["params[date_end]"] = inputParams.DateEnd.Format(dateLayout)
	}
	// sale_id
	addSliceToParams("sale_id", params, inputParams.SaleID)
	// stage
	addSliceToParams("stage", params, inputParams.Stage)

	// Получение ответа

	var resp SalesGetByChangeStageResponse
	if err := request(ctx, apiKey, methodURL, params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
