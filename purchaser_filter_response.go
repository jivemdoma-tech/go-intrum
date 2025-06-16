package gointrum

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type PurchaserFilterResponse struct {
	*Response
	Data *PurchaserFilterData `json:"data,omitempty"`
}

type PurchaserFilterData struct {
	List  []Purchaser `json:"list"`
	Count bool        `json:"count"`
}

type Purchaser struct {
	ID                   uint64                 `json:"id,string"`
	GroupID              string                 `json:"group_id"`
	Name                 string                 `json:"name"`
	Surname              string                 `json:"surname"`
	Secondname           string                 `json:"secondname"`
	ManagerID            uint64                 `json:"manager_id,string"`
	Phone                []*Phone               `json:"phone"`
	Address              string                 `json:"address"`
	CreateDate           time.Time              `json:"create_date"`
	Comment              string                 `json:"comment"`
	CustomerActivityType string                 `json:"customer_activity_type"`
	CustomerActivityDate time.Time              `json:"customer_activity_date"`
	CustomerCreatorID    uint64                 `json:"customer_creator_id,string"`
	Fields               map[uint64]*PurchaserField `json:"fields"`
	EmployeeID           uint64                 `json:"employee_id,string"`
	AdditionalManagerID  []uint64               `json:"additional_manager_id"`
	AdditionalEmployeeID []uint64               `json:"additional_employee_id"`
	// Markname             string                     `json:"markname"`
	// Marktype             string                     `json:"marktype"`
	// Nattype              string                     `json:"nattype"`
	// Email                []interface{}              `json:"email"`

}

type PurchaserField struct {
	Datatype string `json:"datatype"`
	Value    any    `json:"value"`
}

type Phone struct {
	Phone   string `json:"phone"`
	Comment string `json:"comment"`
}

func (p *Purchaser) UnmarshalJSON(data []byte) error {
	// Оригинальная структура типа Alias для предовтращения рекурсии
	type Alias Purchaser

	// Вспомогательная структура
	var aux = &struct {
		*Alias
		CreateDate           string            `json:"create_date"`
		CustomerActivityDate string            `json:"customer_activity_date"`
		AdditionalManagerID  []string          `json:"additional_manager_id"`
		AdditionalEmployeeID []string          `json:"additional_employee_id"`
		Fields               []*PurchaserField `json:"fields"`
	}{
		Alias: (*Alias)(p), // Приведение типа к Alias
	}
	// Декодирование JSON во вспомогательную структуру
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Замена

	parsedDate, err := time.Parse(datetimeLayout, aux.CreateDate)
	if err != nil {
		return err
	}
	p.CreateDate = parsedDate

	parsedDate, err = time.Parse(datetimeLayout, aux.CustomerActivityDate)
	if err != nil {
		return err
	}
	p.CustomerActivityDate = parsedDate

	// Замена массивов

	newSlice := make([]uint64, 0, len(aux.AdditionalManagerID))
	for _, v := range aux.AdditionalManagerID {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	p.AdditionalManagerID = newSlice

	newSlice = make([]uint64, 0, len(aux.AdditionalEmployeeID))
	for _, v := range aux.AdditionalEmployeeID {
		if value, err := strconv.ParseUint(v, 10, 64); err == nil {
			newSlice = append(newSlice, value)
		}
	}
	p.AdditionalEmployeeID = newSlice

	return nil
}

// Методы получения значений Purchaser

// Вспомогательная функция получения структуры поля
func (p *Purchaser) getField(fieldID uint64) (*PurchaserField, bool) {
	f, exists := p.Fields[fieldID]
	return f, exists
}

func (p *Purchaser) getFieldMap(fieldID uint64) (map[string]string, bool) {
	f, exists := p.getField(fieldID)
	if !exists {
		return nil, false
	}
	m, ok := f.Value.(map[string]string)
	if !ok {
		return nil, false
	}
	return m, true
}

// text
func (p *Purchaser) GetFieldText(fieldID uint64) string {
	f, exists := p.getField(fieldID)
	if !exists {
		return ""
	}
	vStr, ok := f.Value.(string)
	if !ok {
		return ""
	}
	return vStr
}

// radio
func (p *Purchaser) GetFieldRadio(fieldID uint64) bool {
	vStr := p.GetFieldText(fieldID)
	if v, err := strconv.ParseBool(vStr); err == nil {
		return v
	}
	return false
}

// select
func (p *Purchaser) GetFieldSelect(fieldID uint64) string {
	return p.GetFieldText(fieldID)
}

// multiselect
func (p *Purchaser) GetFieldMultiselect(fieldID uint64) []string {
	return strings.Split(p.GetFieldText(fieldID), ",")
}

// date
func (p *Purchaser) GetFieldDate(fieldID uint64) time.Time {
	vStr := p.GetFieldText(fieldID)
	return parseTime(vStr, dateLayout)
}

// datetime
func (p *Purchaser) GetFieldDatetime(fieldID uint64) time.Time {
	vStr := p.GetFieldText(fieldID)
	return parseTime(vStr, datetimeLayout)
}

// time
func (p *Purchaser) GetFieldTime(fieldID uint64) time.Time {
	vStr := p.GetFieldText(fieldID)
	return parseTime(vStr, timeLayout)
}

// integer
func (p *Purchaser) GetFieldInteger(fieldID uint64) int64 {
	vStr := p.GetFieldText(fieldID)
	return parseInt(vStr)
}

// decimal
func (p *Purchaser) GetFieldDecimal(fieldID uint64) float64 {
	vStr := p.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// price
func (p *Purchaser) GetFieldPrice(fieldID uint64) float64 {
	vStr := p.GetFieldText(fieldID)
	return parseFloat(vStr)
}

// file
func (p *Purchaser) GetFieldFile(fieldID uint64) string {
	return p.GetFieldText(fieldID)
}

// point
func (p *Purchaser) GetFieldPoint(fieldID uint64) [2]string {
	m, ok := p.getFieldMap(fieldID)
	if !ok {
		return [2]string{}
	}
	return [2]string{m["x"], m["y"]}
}

// integer_range
func (p *Purchaser) GetFieldIntegerRange(fieldID uint64) [2]int64 {
	m, ok := p.getFieldMap(fieldID)
	if !ok {
		return [2]int64{}
	}
	return parseRange(m, parseInt)
}

// decimal_range
func (p *Purchaser) GetFieldDecimalRange(fieldID uint64) [2]float64 {
	m, ok := p.getFieldMap(fieldID)
	if !ok {
		return [2]float64{}
	}
	return parseRange(m, parseFloat)
}

// date_range
func (p *Purchaser) GetFieldDateRange(fieldID uint64) [2]time.Time {
	m, ok := p.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(p string) time.Time {
		return parseTime(p, dateLayout)
	})
}

// time_range
func (p *Purchaser) GetFieldTimeRange(fieldID uint64) [2]time.Time {
	m, ok := p.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(p string) time.Time {
		return parseTime(p, dateLayout)
	})
}

// datetime_range
func (p *Purchaser) GetFieldDatetimeRange(fieldID uint64) [2]time.Time {
	m, ok := p.getFieldMap(fieldID)
	if !ok {
		return [2]time.Time{}
	}
	return parseRange(m, func(p string) time.Time {
		return parseTime(p, dateLayout)
	})
}

// attach
func (p *Purchaser) GetFieldAttach(fieldID uint64) []uint64 {
	f, exists := p.getField(fieldID)
	if !exists {
		return nil
	}
	vAttach, ok := f.Value.([]map[string]string)
	if !ok || len(vAttach) <= 0 {
		return nil
	}
	vIDs := make([]uint64, 0, len(vAttach))
	for _, v := range vAttach {
		if id, err := strconv.ParseUint(v["id"], 10, 64); err == nil {
			vIDs = append(vIDs, id)
		}
	}
	return vIDs
}
