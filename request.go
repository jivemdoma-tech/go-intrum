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
	statusAccessDeny         string = "ACCESS_DENY"
	statusLimitExceeded      string = "LIMIT_EXCEEDED"
	statusBadRequest         string = "BAD_REQUEST"
	statusBadParams          string = "BAD_PARAMS"
	statusServerIsOverloaded string = "SERVER_IS_OVERLOADED"
)

// Клиент для запросов к Intrum API
var client = http.DefaultClient

// Интерфейс, принимающий структуру API-ответа.
type respStruct interface {
	GetErrorMessage() string
}

func request(ctx context.Context, apiKey, reqURL string, reqParams map[string]string, r respStruct) (err error) {
	const (
		primaryPort  string        = "81"
		backupPort   string        = "80"
		duration1Min time.Duration = time.Minute
		duration5Min time.Duration = time.Minute * 5
	)
	// Обработка паники
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occured: %v", r)
		}
	}()

	// URL запроса
	u, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return err
	}
	u.Host = u.Hostname() + ":" + primaryPort // Основной порт

	// Параметры запроса
	p := make(url.Values, len(reqParams)+1)
	p.Set("apikey", apiKey)
	for k, v := range reqParams {
		p.Set(k, v)
	}

	// Запрос

	// Цикл для повторного запроса на запасной порт в случае ошибки
	for _, isBackup := range []bool{false, true} {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if isBackup {
			u.Host = u.Hostname() + ":" + backupPort // Запасной порт
		}

		// Создание нового запроса
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(p.Encode()))
		if err != nil {
			return fmt.Errorf("failed to create request for method %s: %w", u.Path, err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// Отправка запроса
		resp, err := client.Do(req)
		if err != nil {
			if isBackup {
				return fmt.Errorf("failed to do request for method %s: %w", u.Path, err)
			}
			// Повторный запрос
			time.Sleep(duration1Min)
			continue
		}

		// Ответ

		// Чтение ответа
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Закрытие чтения тела ответа
		if err != nil {
			return fmt.Errorf("failed to read response body from method %s: %w", u.Path, err)
		}
		// Non-2xx status code
		if resp.StatusCode >= http.StatusMultipleChoices && resp.StatusCode != http.StatusNotImplemented {
			if isBackup {
				return fmt.Errorf("%d status code from method %s", resp.StatusCode, u.Path)
			}
			// Повторный запрос
			switch resp.StatusCode {
			case http.StatusGatewayTimeout:
				time.Sleep(duration5Min)
				continue
			default:
				time.Sleep(duration1Min)
				continue
			}
		}
		// Декодирование ответа
		if err := json.Unmarshal(body, r); err != nil {
			return fmt.Errorf("failed to decode response body from method %s: %w", u.Path, err)
		}
		// Обработка ошибки от сервера
		if errMsg := r.GetErrorMessage(); errMsg != "" {
			if isBackup {
				return fmt.Errorf("response error from method %s: %s", u.Path, errMsg)
			}
			switch {
			case
				strings.Contains(errMsg, statusAccessDeny),    // ACCESS_DENY
				strings.Contains(errMsg, statusLimitExceeded), // LIMIT_EXCEEDED
				strings.Contains(errMsg, statusBadRequest),    // BAD_REQUEST
				strings.Contains(errMsg, statusBadParams):     // BAD_PARAMS
				return fmt.Errorf("response error from method %s: %s", u.Path, errMsg)
			// Повторный запрос через 5 минут
			case strings.Contains(errMsg, statusServerIsOverloaded):
				time.Sleep(duration5Min)
				continue
			// Повторный запрос через минуту
			default:
				time.Sleep(duration1Min)
				continue
			}
		}

		break
	}

	return nil
}
