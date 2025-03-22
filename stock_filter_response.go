package gointrum

import "time"

type StockFilterResponse struct {
	Status string           `json:"status"`
	Data   *StockFilterData `json:"data"`
}
type StockFilterData struct {
	List []*StockFilterList `json:"list"`
	// Count bool               `json:"count"`
}
type StockFilterList struct {
	ID                   uint64              `json:"id,string"`
	StockType            uint16              `json:"stock_type,string"`
	Type                 uint16              `json:"type,string"`
	Parent               uint16              `json:"parent,string"`
	Name                 string              `json:"name"`
	DateAdd              time.Time           `json:"date_add"` // TODO
	Count                bool                `json:"count,string"`
	Author               uint64              `json:"author,string"`
	EmployeeID           uint64              `json:"employee_id,string"`
	AdditionalAuthor     []uint64            `json:"additional_author"`
	AdditionalEmployeeID []uint64            `json:"additional_employee_id"`
	LastModify           time.Time           `json:"last_modify"` // TODO
	CustomerRelation     uint64              `json:"customer_relation,string"`
	StockActivityType    string              `json:"stock_activity_type"`
	StockActivityDate    time.Time           `json:"stock_activity_date"` // TODO
	Publish              bool                `json:"publish,string"`
	Copy                 uint64              `json:"copy,string"`
	GroupID              uint16              `json:"group_id,string"`
	StockCreatorID       uint64              `json:"stock_creator_id,string"`
	Fields               []*StockFilterField `json:"fields"`
	// Log                  interface{}       `json:"log"`
}
type StockFilterField struct {
	ID    uint64 `json:"id,string"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}
