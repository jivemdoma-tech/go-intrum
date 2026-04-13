package intrum

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
	primaryPort      string = "81"
	backupPort       string = "444"
	backupDelayShort        = 30 * time.Second
	backupDelayLong         = 5 * time.Minute

	statusAccessDeny         string = "ACCESS_DENY"
	statusLimitExceeded      string = "LIMIT_EXCEEDED"
	statusBadRequest         string = "BAD_REQUEST"
	statusBadParams          string = "BAD_PARAMS"
	statusServerIsOverloaded string = "SERVER_IS_OVERLOADED"
)

// Клиент для запросов к Intrum API
var client = http.Client{Timeout: 10 * time.Minute}

// Интерфейс, принимающий структуру API-ответа.
type respStruct interface {
	GetErrorMessage() string
}

func request(ctx context.Context, apiKey, reqURL string, reqParams map[string]string, r respStruct) (err error) {
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
		if isBackup {
			var (
				hostNameParts = strings.Split(u.Hostname(), ".")
				hostName      = "intrum4." + strings.Join(hostNameParts[1:], ".")
			)
			u.Host = hostName + ":" + backupPort // Запасной порт
			u.Scheme = "https"
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
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backupDelayShort):
				continue
			}
		}

		// Ответ

		// Чтение ответа
		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close() // Закрытие чтения тела ответа
		if err != nil {
			return fmt.Errorf("failed to read response body from method %s: %w", u.Path, err)
		}
		switch {
		// Non-2xx status code
		case resp.StatusCode >= http.StatusMultipleChoices && resp.StatusCode != http.StatusNotImplemented:
			if isBackup {
				return fmt.Errorf("%d status code from method %s", resp.StatusCode, u.Path)
			}
			// Таймаут + повторный запрос
			timeout := func() time.Duration {
				switch resp.StatusCode {
				// Ошибка на стороне сервера (502, 504)
				case http.StatusBadGateway, http.StatusGatewayTimeout:
					return backupDelayLong
				// Другая ошибка
				default:
					return backupDelayShort
				}
			}()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(timeout):
				continue
			}
		// Пустое тело ответа
		case len(body) == 0:
			if isBackup {
				return fmt.Errorf("empty response body with status code %d from method %s", resp.StatusCode, u.Path)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backupDelayShort):
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
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backupDelayLong):
					continue
				}
			// Повторный запрос через минуту
			default:
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backupDelayShort):
					continue
				}
			}
		}

		break
	}

	return nil
}
