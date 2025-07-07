package gointrum

import (
	"encoding/json"
	"strconv"
)

type SalesDetailsResponse struct {
	*Response
	Data map[int64]*SalesDetailsData `json:"data,omitempty"`
}

type Data struct {
	Sales SalesDetailsData `json:"sales"`
}

type SalesDetailsData struct {
	Cocustomers    []int64 `json:"cocustomers"`
	Corequests     []int64 `json:"corequests"`
	Costock        []int64 `json:"costock"`
	Costockprimary []int64 `json:"costockprimary"`
	Bills          []int64 `json:"bills"`
	Blanks         []int64 `json:"blanks"`
	Cosales        []int64 `json:"cosales"`
}

func (s *SalesDetailsResponse) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias SalesDetailsResponse

	// Вспомогательная структура
	aux := &struct {
		*Alias
		Data map[string]*SalesDetailsData `json:"data,omitempty"`
	}{
		Alias: (*Alias)(s),
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена массивов

	s.Data = func() map[int64]*SalesDetailsData {
		newMap := make(map[int64]*SalesDetailsData, len(aux.Data))
		for k, v := range aux.Data {
			kInt, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}
			newMap[kInt] = v
		}
		switch len(newMap) {
		case 0:
			return nil
		default:
			return newMap
		}
	}()

	return nil
}

func (s *SalesDetailsData) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias SalesDetailsData

	// Вспомогательная структура
	aux := &struct {
		*Alias
		Cocustomers    []string `json:"cocustomers"`
		Corequests     []string `json:"corequests"`
		Costock        []string `json:"costock"`
		Costockprimary []string `json:"costockprimary"`
		Bills          []string `json:"bills"`
		Blanks         []string `json:"blanks"`
		Cosales        []string `json:"cosales"`
	}{
		Alias: (*Alias)(s),
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена массивов

	s.Cocustomers = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Cocustomers))
		for _, v := range aux.Cocustomers {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	s.Corequests = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Corequests))
		for _, v := range aux.Corequests {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	s.Costock = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Costock))
		for _, v := range aux.Costock {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	s.Costockprimary = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Costockprimary))
		for _, v := range aux.Costockprimary {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	s.Bills = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Bills))
		for _, v := range aux.Bills {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	s.Blanks = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Blanks))
		for _, v := range aux.Blanks {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	s.Cosales = func() []int64 {
		newSlice := make([]int64, 0, len(aux.Cosales))
		for _, v := range aux.Cosales {
			vInt, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				continue
			}
			newSlice = append(newSlice, vInt)
		}
		switch len(newSlice) {
		case 0:
			return nil
		default:
			return newSlice
		}
	}()

	return nil
}
