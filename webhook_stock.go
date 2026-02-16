package gointrum

import (
	"fmt"
	"strconv"
	"strings"
)

type WebhookStockPayload struct {
	SubjectType   string               `json:"subject_type"`
	SubjectTypeID int64                `json:"subject_type_id,string"`
	Event         string               `json:"event"`
	ObjectType    string               `json:"object_type"`
	ObjectSubType int64                `json:"object_sub_type"`
	ObjectSubID   int64                `json:"object_sub_id"`
	Snapshot      WebhookStockSnapshot `json:"snapshot"`
}

type WebhookStockSnapshot struct {
	ID                  int64             `json:"id,string"`
	GroupID             string            `json:"group_id"`
	ExportID            string            `json:"export_id"`
	Parent              string            `json:"parent"`
	Parentname          string            `json:"parentname"`
	ExportParent        string            `json:"export_parent"`
	Name                string            `json:"name"`
	Count               string            `json:"count"`
	DateAdd             string            `json:"date_add"`
	Publish             string            `json:"publish"`
	PostProcessed       string            `json:"post_processed"`
	Author              string            `json:"author"`
	Copy                string            `json:"copy"`
	Type                string            `json:"type"`
	Typename            string            `json:"typename"`
	StockActivityType   string            `json:"stock_activity_type"`
	StockActivityDate   string            `json:"stock_activity_date"`
	RelatedWithCustomer string            `json:"related_with_customer"`
	StockCreatorId      string            `json:"stock_creator_id"`
	EntitySubId         string            `json:"entity_sub_id"`
	MainOwnerId         string            `json:"main_owner_id"`
	SharedManagers      []any             `json:"shared_managers"`
	Extproperty         WebhookStockExtpr `json:"extproperty"`
	Merge               WebhookStockMerge `json:"merge"`
}

type (
	WebhookStockExtpr      map[string]*WebhookStockExtprField
	WebhookStockExtprField struct {
		Type  string `json:"type"`
		Value any    `json:"value"`
	}
)

func (ff WebhookStockExtpr) field(fieldID int64) *WebhookStockExtprField {
	if len(ff) == 0 {
		return nil
	}

	f, ok := ff[strconv.FormatInt(fieldID, 10)]
	if !ok {
		return nil
	}

	return f
}

func (ff WebhookStockExtpr) StringField(fieldID int64) string {
	f := ff.field(fieldID)
	if f == nil {
		return ""
	}
	v := f.Value
	if v == nil {
		return ""
	}

	return fmt.Sprint(v)
}

func (ff WebhookStockExtpr) BoolField(fieldID int64) bool {
	vStr := ff.StringField(fieldID)
	if vStr == "" {
		return false
	}

	vBool, _ := strconv.ParseBool(vStr)
	return vBool
}

func (ff WebhookStockExtpr) FloatField(fieldID int64) float64 {
	vStr := ff.StringField(fieldID)
	if vStr == "" {
		return 0
	}

	vFloat64, _ := strconv.ParseFloat(vStr, 64)
	return vFloat64
}

func (ff WebhookStockExtpr) IntField(fieldID int64) int64 {
	vFloat := ff.FloatField(fieldID)
	return int64(vFloat)
}

func (ff WebhookStockExtpr) StringSliceField(fieldID int64) []string {
	vStr := ff.StringField(fieldID)
	if vStr == "" {
		return nil
	}

	vStrSlice := strings.Split(vStr, ",")
	if len(vStrSlice) == 0 {
		return nil
	}

	return vStrSlice
}

type (
	WebhookStockMerge      map[string]*WebhookStockMergeField
	WebhookStockMergeField struct {
		Type    string `json:"type"`
		Value   any    `json:"value"`
		Current any    `json:"current"`
	}
	WebhookStockMergeEdit[T any] struct {
		Before T
		After  T
	}
)

func (ff WebhookStockMerge) fieldEdit(fieldID int64) *WebhookStockMergeField {
	if len(ff) == 0 {
		return nil
	}

	f, ok := ff[strconv.FormatInt(fieldID, 10)]
	if !ok {
		return nil
	}

	return f
}

func (ff WebhookStockMerge) StringFieldEdit(fieldID int64) *WebhookStockMergeEdit[string] {
	f := ff.fieldEdit(fieldID)
	if f == nil {
		return nil
	}

	var vBefore string
	if v := f.Value; v != nil {
		vBefore = fmt.Sprint(v)
	}

	var vAfter string
	if v := f.Current; v != nil {
		vAfter = fmt.Sprint(v)
	}

	return &WebhookStockMergeEdit[string]{
		Before: vBefore,
		After:  vAfter,
	}
}

func (ff WebhookStockMerge) BoolFieldEdit(fieldID int64) *WebhookStockMergeEdit[bool] {
	editStr := ff.StringFieldEdit(fieldID)
	if editStr == nil {
		return nil
	}

	var (
		vBeforeBool, _ = strconv.ParseBool(editStr.Before)
		vAfterBool, _  = strconv.ParseBool(editStr.After)
	)

	return &WebhookStockMergeEdit[bool]{
		Before: vBeforeBool,
		After:  vAfterBool,
	}
}

func (ff WebhookStockMerge) FloatFieldEdit(fieldID int64) *WebhookStockMergeEdit[float64] {
	editStr := ff.StringFieldEdit(fieldID)
	if editStr == nil {
		return nil
	}

	var (
		vBeforeFloat64, _ = strconv.ParseFloat(editStr.Before, 64)
		vAfterFloat64, _  = strconv.ParseFloat(editStr.After, 64)
	)

	return &WebhookStockMergeEdit[float64]{
		Before: vBeforeFloat64,
		After:  vAfterFloat64,
	}
}

func (ff WebhookStockMerge) IntFieldEdit(fieldID int64) *WebhookStockMergeEdit[int64] {
	editFloat := ff.FloatFieldEdit(fieldID)
	if editFloat == nil {
		return nil
	}

	return &WebhookStockMergeEdit[int64]{
		Before: int64(editFloat.Before),
		After:  int64(editFloat.After),
	}
}

func (ff WebhookStockMerge) StringSliceFieldEdit(fieldID int64) *WebhookStockMergeEdit[[]string] {
	editStr := ff.StringFieldEdit(fieldID)
	if editStr == nil {
		return nil
	}

	vBeforeStrSlice := strings.Split(editStr.Before, ",")
	if len(vBeforeStrSlice) == 0 {
		vBeforeStrSlice = nil
	}

	vAfterStrSlice := strings.Split(editStr.After, ",")
	if len(vAfterStrSlice) == 0 {
		vAfterStrSlice = nil
	}

	return &WebhookStockMergeEdit[[]string]{
		Before: vBeforeStrSlice,
		After:  vAfterStrSlice,
	}
}
