package gointrum

// import "encoding/json"

// type TopLevel struct {
// 	*Response
// 	Data   map[string]*StockAttachData `json:"data,omitempty"`
// }

// type Data struct {
// 	The41006 The41006 `json:"41006"`
// }

// type The41006 struct {
// 	Cocustomers    []string      `json:"cocustomers"`
// 	Corequests     []interface{} `json:"corequests"`
// 	Costock        []string      `json:"costock"`
// 	Costockprimary []string      `json:"costockprimary"`
// 	Bills          []interface{} `json:"bills"`
// 	Blanks         []string      `json:"blanks"`
// 	Cosales        []interface{} `json:"cosales"`
// }

// type StockAttachResponse struct {
// 	*Response
// 	Data map[string]*StockAttachData `json:"data,omitempty"`
// }

// type StockAttachData struct {
// 	Requests []string `json:"requests"`
// }

// func (s *StockAttachResponse) UnmarshalJSON(data []byte) error {
// 	type Alias StockAttachResponse

// 	// Временная структура с raw json для поля Data
// 	aux := &struct {
// 		Data json.RawMessage `json:"data"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(s),
// 	}

// 	// Парсим все кроме поля Data
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}

// 	// Обработка поля Data
// 	if string(aux.Data) == "[]" {
// 		s.Data = nil
// 	} else {
// 		// обычный парсинг в map
// 		if err := json.Unmarshal(aux.Data, &s.Data); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
