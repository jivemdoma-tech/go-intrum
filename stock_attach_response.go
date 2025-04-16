package gointrum

import "encoding/json"

type StockAttachResponse struct {
	*Response
	Data map[string]*StockAttachData `json:"data,omitempty"`
}

type StockAttachData struct {
	Requests []string `json:"requests"`
}

func (s *StockAttachResponse) UnmarshalJSON(data []byte) error {
	type Alias StockAttachResponse

	// Временная структура с raw json для поля Data
	aux := &struct {
		Data json.RawMessage `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	// Парсим все кроме поля Data
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Обработка поля Data
	if string(aux.Data) == "[]" {
		s.Data = nil
	} else {
		// обычный парсинг в map
		if err := json.Unmarshal(aux.Data, &s.Data); err != nil {
			return err
		}
	}

	return nil
}
