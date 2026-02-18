package gointrum

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type WorkerFilterResponse struct {
	*Response
	Data map[string]*WorkerFilterData `json:"data"`
}

type WorkerFilterData struct {
	ID          string                   `json:"id"`
	Type        string                   `json:"type"`
	DivisionID  string                   `json:"division_id"`
	SubofficeID string                   `json:"suboffice_id"`
	Post        string                   `json:"post"`
	Boss        string                   `json:"boss"`
	Status      string                   `json:"status"`
	Name        string                   `json:"name"`
	Surname     string                   `json:"surname"`
	Fields      map[uint64]*workerFields `json:"fields"`
	// Secondname          string           `json:"secondname"`
	// Internalemail       []string         `json:"internalemail"`
	// Externalemail       []interface{}    `json:"externalemail"`
	// Internalphone       []interface{}    `json:"internalphone"`
	// Externalphone       []interface{}    `json:"externalphone"`
	// Mobilephone         []Mobilephone    `json:"mobilephone"`
	// Birthday            *string          `json:"birthday"`
	// Address             interface{}      `json:"address"`
	// About               string           `json:"about"`
	// Hobby               string           `json:"hobby"`
	// CreatedAt           *time.Time       `json:"created_at"`
	// Skype               string           `json:"skype"`
	// Facebook            *string          `json:"facebook"`
	// Vkontakte           string           `json:"vkontakte"`
	// Gender              string           `json:"gender"`
	// Avatars             Avatars          `json:"avatars"`
	// AsteriskShortNumber []string         `json:"asterisk_short_number,omitempty"`
}

type workerFields struct {
	ID       uint64 `json:"id,string,omitempty"`
	Name     any    `json:"value,omitempty"`
	Datatype string `json:"datatype,omitempty"`
}

func (w *WorkerFilterData) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias WorkerFilterData

	// Вспомогательная структура
	aux := &struct {
		*Alias
		Fields any `json:"fields"`
	}{
		Alias: (*Alias)(w),
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if m, ok := aux.Fields.(map[string]any); ok {
		out := make(map[uint64]*workerFields, len(m))
		for k, v := range m {
			id, err := strconv.ParseUint(k, 10, 64)
			if err != nil {
				// пропускаем некорректный ключ
				continue
			}
			b, _ := json.Marshal(v)
			var f workerFields
			if err := json.Unmarshal(b, &f); err != nil {
				continue
			}
			// если в объекте нет id, всё равно ключ у нас верный — установим
			if f.ID == 0 {
				f.ID, _ = strconv.ParseUint(k, 10, 64)
			}
			out[id] = &f
		}
		w.Fields = out
		return nil
	}
	return nil
}

// Методы получения значений полей

// getField получает структуру поля по ID.
func (w *WorkerFilterData) getField(fieldID uint64) *workerFields {
	if f, exists := w.Fields[fieldID]; exists {
		return f
	}
	return nil
}

func (w *WorkerFilterData) getFieldMap(fieldID uint64) map[string]string {
	f := w.getField(fieldID)
	if f == nil {
		return nil
	}
	switch m := f.Name.(type) {
	case map[string]string:
		return m
	case map[string]any:
		mStr := make(map[string]string, len(m))
		for k, v := range m {
			mStr[k] = fmt.Sprint(v)
		}
		return mStr
	}
	return nil
}

// Публичные методы

// GetFieldText возвращает string значение поля.
func (w *WorkerFilterData) GetFieldText(fieldID uint64) string {
	f := w.getField(fieldID)
	if f == nil {
		return ""
	}
	vStr, ok := f.Name.(string)
	if !ok {
		return ""
	}
	return vStr
}

// GetFieldSelect возвращает string значение поля.
func (w *WorkerFilterData) GetFieldSelect(fieldID uint64) string {
	return w.GetFieldText(fieldID)
}

// GetFieldInteger возвращает int64 значение поля.
func (w *WorkerFilterData) GetFieldInteger(fieldID uint64) int64 {
	vStr := w.GetFieldText(fieldID)
	return parseInt(vStr)
}

// GetFieldDecimal возвращает float64 значение поля.
func (w *WorkerFilterData) GetFieldDecimal(fieldID uint64) float64 {
	vStr := w.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// TODO все остальные типы данных
