package gointrum

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UploadFile struct {
	FieldName string
	FileName  string
	R         io.Reader
	Closer    io.Closer
}

func requestUploadFile(ctx context.Context, apiKey, reqURL string, reqParams map[string]string, files []UploadFile, r respStruct) (err error) {
	const (
		primaryPort  string = "81"
		backupPort   string = "444"
		duration1Min        = time.Minute
		duration5Min        = time.Minute * 5
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
			var (
				hostNameParts = strings.Split(u.Hostname(), ".")
				hostName      = "intrum4." + strings.Join(hostNameParts[1:], ".")
			)
			u.Host = hostName + ":" + backupPort // Запасной порт
			u.Scheme = "https"
		}

		var (
			errCh       chan error // будет nil, если не multipart
			body        io.Reader
			contentType string
		)

		if len(files) > 0 {
			pr, pw := io.Pipe()
			w := multipart.NewWriter(pw)

			contentType = w.FormDataContentType()
			body = pr
			errCh = make(chan error, 1)

			go func() {
				// ВАЖНО: любая ошибка -> CloseWithError, чтобы клиент не висел
				fail := func(e error) {
					_ = w.Close()
					_ = pw.CloseWithError(e)
					errCh <- e
				}

				// поля
				if e := w.WriteField("apikey", apiKey); e != nil {
					fail(e)
					return
				}
				for k, v := range reqParams {
					if e := w.WriteField(k, v); e != nil {
						fail(e)
						return
					}
				}

				// файлы (стриминг)
				for _, f := range files {
					part, e := w.CreateFormFile(f.FieldName, f.FileName)
					if e != nil {
						fail(e)
						return
					}
					if _, e = io.Copy(part, f.R); e != nil {
						fail(e)
						return
					}
					if f.Closer != nil {
						_ = f.Closer.Close()
					}
				}

				// закрываем multipart и pipe
				if e := w.Close(); e != nil {
					_ = pw.CloseWithError(e)
					errCh <- e
					return
				}
				_ = pw.Close()
				errCh <- nil
			}()
		} else {
			body = strings.NewReader(p.Encode())
			contentType = "application/x-www-form-urlencoded"
		}

		// Создание нового запроса
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), body)
		if err != nil {
			return fmt.Errorf("failed to create request for method %s: %w", u.Path, err)
		}
		req.Header.Set("Content-Type", contentType)
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

		if errCh != nil {
			if e := <-errCh; e != nil {
				resp.Body.Close()
				return fmt.Errorf("multipart build failed: %w", e)
			}
		}

		// Ответ

		// Чтение ответа
		bodyResp, err := io.ReadAll(resp.Body)
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
		if err := json.Unmarshal(bodyResp, r); err != nil {
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
