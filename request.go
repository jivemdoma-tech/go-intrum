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

// Клиент для запросов к Intrum API
var client = &http.Client{
	Timeout: time.Duration(10 * time.Minute),
}

// Интерфейс структуры API-ответа
type respStruct interface {
	GetErrorMessage() string
}

func rawRequest(ctx context.Context, apiKey, u string, p map[string]string, r respStruct) error {
	params := make(url.Values, len(p)+1)
	params.Set("apikey", apiKey) // Параметр, содержащий API-ключ
	for k, v := range p {
		params.Set(k, v)
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
		uBackup := strings.Replace(u, "81", "80", -1)
		// Запасной запрос
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, uBackup, strings.NewReader(params.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create request for method %s: %w", uBackup, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// Отправка запасного запроса на сервер
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to do request for method %s: %w", uBackup, err)
		}
	}
	defer resp.Body.Close()

	// Обработка ответа от сервера

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body from method %s: %w", u, err)
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("status code %d from method %s", resp.StatusCode, u)
	}

	// Декодирование ответа

	if err := json.Unmarshal(body, r); err != nil {
		return fmt.Errorf("failed to decode response body from method %s: %w", u, err)
	}

	// Проверка ответ с ошибкой // TODO: Разгрузить логику

	switch r.GetErrorMessage() {
	case "":
		break
	default:
		uBackup := strings.Replace(u, "81", "80", -1)
		// Запасной запрос
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, uBackup, strings.NewReader(params.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create request for method %s: %w", uBackup, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// Отправка запасного запроса на сервер
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to do request for method %s: %w", uBackup, err)
		}
		defer resp.Body.Close()
		// Обработка ответа от сервера
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body from method %s: %w", uBackup, err)
		}
		if resp.StatusCode >= http.StatusMultipleChoices {
			return fmt.Errorf("status code %d from method %s", resp.StatusCode, uBackup)
		}
		// Декодирование ответа
		if err := json.Unmarshal(body, r); err != nil {
			return fmt.Errorf("failed to decode response body from method %s: %w", uBackup, err)
		}
		// Проверка ответ с ошибкой
		if r.GetErrorMessage() != "" {
			return fmt.Errorf("error code %s from method %s", r.GetErrorMessage(), u)
		}
	}

	return nil
}
