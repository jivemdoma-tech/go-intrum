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

const (
	RespStatusServerIsOverloaded string = "SERVER_IS_OVERLOADED"
	RespStatusAccessDeny         string = "ACCESS_DENY"
)

// Клиент для запросов к Intrum API
var client = http.DefaultClient

// Интерфейс структуры API-ответа
type respStruct interface {
	GetErrorMessage() string
}

func rawRequest(ctx context.Context, apiKey, u string, p map[string]string, r respStruct) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occured: %v", r)
		}
	}()

	// Параметры запроса

	params := make(url.Values, len(p)+1)
	params.Set("apikey", apiKey)
	for k, v := range p {
		params.Set(k, v)
	}

	// Запрос

	// Цикл для повторного запроса на запасной порт в случае ошибки
	for _, isBackupRequest := range []bool{false, true} {
		if isBackupRequest {
			u = strings.Replace(u, "81", "80", -1)
		}

		// Формирование нового запроса
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(params.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create request for method %s: %w", u, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// Отправка запроса на сервер
		resp, err := client.Do(req)
		if err != nil {
			if isBackupRequest {
				return fmt.Errorf("failed to do request for method %s: %w", u, err)
			}
			// Запрос на запасной порт
			time.Sleep(time.Minute)
			continue
		}
		defer resp.Body.Close()

		// Ответ

		// Чтение ответа от сервера
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body from method %s: %w", u, err)
		}
		// Non-2xx status code
		if resp.StatusCode >= http.StatusMultipleChoices {
			if isBackupRequest {
				return fmt.Errorf("%d status code from method %s", resp.StatusCode, u)
			}
			// Запрос на запасной порт
			time.Sleep(time.Minute)
			continue
		}
		// Декодирование ответа
		if err := json.Unmarshal(body, r); err != nil {
			return fmt.Errorf("failed to decode response body from method %s: %w", u, err)
		}
		// Повторный запрос при ошибке от сервера
		if errMessage := r.GetErrorMessage(); errMessage != "" {
			if isBackupRequest {
				return fmt.Errorf("error response from method %s: %s", u, errMessage)
			}
			// Запрос на запасной порт
			switch {
			case strings.Contains(errMessage, RespStatusAccessDeny):
				return fmt.Errorf("error response from method %s: %s", u, errMessage)
			case strings.Contains(errMessage, RespStatusServerIsOverloaded):
				time.Sleep(time.Minute * 5)
			default:
				time.Sleep(time.Minute)
			}
			continue
		}
		break
	}
	return nil
}
