package intrumgo

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

const contentType string = "application/x-www-form-urlencoded"

var (
	stdTimeout = time.Duration(time.Second * 300)
	client  = &http.Client{Timeout: stdTimeout}
)

func rawRequest(ctx context.Context, methodURL, apiKey string, timeoutSec int, params map[string]string, r any) error {
	var timeout time.Duration
	switch timeoutSec {
	case 0:
		timeout = stdTimeout
	default:
		timeout = time.Duration(time.Second * 300)
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Параметры запроса

	p := make(url.Values, len(params)+1)
	p.Set("apikey", apiKey)
	for k, v := range params {
		p.Set(k, v)
	}
	httpBody := strings.NewReader(p.Encode())

	// Новый запрос

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, methodURL, httpBody)
	if err != nil {
		return fmt.Errorf("error create request for method %s: %w", methodURL, err)
	}
	req.Header.Set("Content-Type", contentType)

	// Отправка запроса на сервер

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error do request for method %s: %w", methodURL, err)
	}
	defer resp.Body.Close()

	// Обработка ответа от сервера

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response body for method %s: %w", methodURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from intrum for method %s: %d", methodURL, resp.StatusCode)
	}

	// Декодирование ответа

	err = json.Unmarshal(body, r)
	if err != nil {
		return fmt.Errorf("error decode response body for method %s: %w", methodURL, err)
	}

	return nil
}
