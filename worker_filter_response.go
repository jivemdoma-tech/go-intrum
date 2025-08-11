package gointrum

import (
	"encoding/json"
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
	Fields      map[uint64]*WorkerFields `json:"fields"`
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
	// GroupID             []string         `json:"group_id"`
	// Avatars             Avatars          `json:"avatars"`
	// AsteriskShortNumber []string         `json:"asterisk_short_number,omitempty"`
}

type WorkerFields struct {
	ID       uint64 `json:"id,string,omitempty"`
	Name     string `json:"value,omitempty"`
	Datatype any    `json:"datatype,omitempty"`
}

func (w *WorkerFilterData) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias WorkerFilterData

	// Вспомогательная структура
	var aux = &struct {
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
		out := make(map[uint64]*WorkerFields, len(m))
		for k, v := range m {
			id, err := strconv.ParseUint(k, 10, 64)
			if err != nil {
				// пропускаем некорректный ключ
				continue
			}
			b, _ := json.Marshal(v)
			var f WorkerFields
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
