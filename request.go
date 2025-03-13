package gointrum

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Стандартное время ожидания ответа от Intrum API
const stdTimeout time.Duration = time.Duration(10 * time.Minute)

// Клиент для запросов к Intrum API
var client = &http.Client{
	Timeout: stdTimeout,
}

type respStruct interface {
	// Объекты
	*StockInsertResponse |
		// Сделки
		*SalesTypesResponse | *SalesGetByChangeStageResponse |
		*SalesFilterResponse | *SalesUpdateResponse
}

func rawRequest[T respStruct](ctx context.Context, apiKey, u string, p map[string]string, r T) error {
	ctx, cancel := context.WithTimeout(ctx, stdTimeout)
	defer cancel()

	// Параметры запроса

	params := make(url.Values, len(p)+1)
	params.Set("apikey", apiKey)
	for k, v := range p {
		params.Set(k, v)
	}
	httpBody := strings.NewReader(params.Encode())

	// Новый запрос

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, httpBody)
	if err != nil {
		return fmt.Errorf("error create request for method %s: %w", u, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Отправка запроса на сервер

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error do request for method %s: %w", u, err)
	}
	defer resp.Body.Close()

	// Обработка ответа от сервера

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response body for method %s: %w", u, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from intrum for method %s: %d", u, resp.StatusCode)
	}

	// Декодирование ответа

	err = json.Unmarshal(body, r)
	if err != nil {
		return fmt.Errorf("error decode response body for method %s: %w", u, err)
	}

	// TODO: Добавить запрос на альтернативный порт 80 при определенных ответах от сервера

	return nil
}
