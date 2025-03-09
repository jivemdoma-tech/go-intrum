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
	Stage     []uint16
}

// Ссылка на метод: https://www.intrumnet.com/api/#sales-filter-stage-period
func SalesGetByChangeStage(ctx context.Context, subdomain, apiKey string, timeoutSec int, inputParams *SalesGetByChangeStageParams) (*SalesGetByChangeStageResponse, error) {
	var (
		primaryURL string = fmt.Sprintf("http://%s.intrumnet.com:81/sharedapi/sales/getbychangestage", subdomain)
		backupURL  string = fmt.Sprintf("http://%s.intrumnet.com:80/sharedapi/sales/getbychangestage", subdomain)
	)

	// Параметры запроса

	params := make(map[string]string, getParamsSize(inputParams))

	// date_start
	params["params[date_start]"] = inputParams.DateStart.Format(dateLayout)
	// date_end
	params["params[date_end]"] = inputParams.DateEnd.Format(dateLayout)
	// sale_id
	addSliceToParams(params, "sale_id", inputParams.SaleID)
	// stage
	addSliceToParams(params, "stage", inputParams.Stage)

	// Получение ответа

	var resp SalesGetByChangeStageResponse

	if err := rawRequest(ctx, primaryURL, apiKey, timeoutSec, params, &resp); err != nil {
		if err := rawRequest(ctx, backupURL, apiKey, timeoutSec, params, &resp); err != nil {
			return nil, err
		}
	}

	return &resp, nil
}
