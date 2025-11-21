package gointrum

import "encoding/json"

type PurchaserAttachResponse struct {
	*Response
	Data map[string]*PurchaserAttachData `json:"data,omitempty"`
}

type PurchaserAttachData struct {
	Stock           []interface{}     `json:"stock,omitempty"`
	StockExtended   []interface{}     `json:"stock_extended,omitempty"`
	StockArchive    []interface{}     `json:"stock_archive,omitempty"`
	Request         []string          `json:"request,omitempty"`
	RequestExtended []RequestExtended `json:"request_extended,omitempty"`
	Sale            []interface{}     `json:"sale,omitempty"`
	SaleExtended    []interface{}     `json:"sale_extended,omitempty"`
}

type RequestExtended struct {
	ID       string `json:"id"`
	TypeID   string `json:"type_id"`
	TypeName string `json:"type_name"`
}

func (s *PurchaserAttachResponse) UnmarshalJSON(data []byte) error {
	type Alias PurchaserAttachResponse

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
