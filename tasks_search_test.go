package gointrum

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTasksSearch(t *testing.T) {
	// Сохранение оригинального raw request
	ogRequestFn := requestFn
	defer func() { requestFn = ogRequestFn }()
	// Подмена raw request
	requestFn = func(ctx context.Context, apiKey, reqURL string, reqParams map[string]string, r respStruct) (err error) {
		// Неуспешный ответ
		pTitle, ok := reqParams["params[title]"]
		if !ok || strings.TrimSpace(pTitle) == "" {
			return fmt.Errorf("response error: %s", statusBadParams)
		}

		// Успешный ответ
		respJSON, err := os.ReadFile(filepath.Join(".", "tasks_search_test.json"))
		assertNoErr(t, err)
		err = json.Unmarshal(respJSON, r)
		assertNoErr(t, err)

		return nil
	}

	t.Run("Request with bad params", func(t *testing.T) {
		_, err := TasksSearch(context.Background(), "test-subdomain", "test-api-key", TasksSearchParams{})
		assertErr(t, err)
	})

	t.Run("Successful request", func(t *testing.T) {
		resp, err := TasksSearch(context.Background(), "test-subdomain", "test-api-key", TasksSearchParams{Title: "тест"})
		assertNoErr(t, err)

		// Проверка что хоть в одной итерации по JSON все поля обработаны верно

		var (
			checkID          bool
			checkCreatedAt   bool
			checkTitle       bool
			checkDescription bool
			checkStatus      bool
			checkPriority    bool
			checkDirector    bool
			checkPerformer   bool
			checkCoperformer bool
			checkAuthor      bool
			checkAttaches    bool
		)

		for _, task := range resp.Data.Tasks {
			if task.ID != 0 {
				checkID = true
			}
			if !task.CreatedAt.IsZero() {
				checkCreatedAt = true
			}
			if task.Title != "" {
				checkTitle = true
			}
			if task.Description != "" {
				checkDescription = true
			}
			if task.Status != "" {
				checkStatus = true
			}
			if task.Priority != 0 {
				checkPriority = true
			}
			if task.Author != 0 {
				checkAuthor = true
			}
			if task.Director != 0 {
				checkDirector = true
			}
			if task.Performer != 0 {
				checkPerformer = true
			}
			if len(task.Coperformer) != 0 {
				checkCoperformer = true
			}
			if len(task.Attaches) != 0 {
				checkAttaches = true
			}
		}

		switch {
		case !checkID:
			assertNoErr(t, errors.New("missing id"))
		case !checkCreatedAt:
			assertNoErr(t, errors.New("missing created_at"))
		case !checkTitle:
			assertNoErr(t, errors.New("missing title"))
		case !checkDescription:
			assertNoErr(t, errors.New("missing description"))
		case !checkStatus:
			assertNoErr(t, errors.New("missing status"))
		case !checkPriority:
			assertNoErr(t, errors.New("missing priority"))
		case !checkDirector:
			assertNoErr(t, errors.New("missing director"))
		case !checkPerformer:
			assertNoErr(t, errors.New("missing performer"))
		case !checkCoperformer:
			assertNoErr(t, errors.New("missing coperformer"))
		case !checkAuthor:
			assertNoErr(t, errors.New("missing author"))
		case !checkAttaches:
			assertNoErr(t, errors.New("missing attaches"))
		}
	})
}
