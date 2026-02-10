package gointrum

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	SubjectTypeSystem   string = "system"
	SubjectTypeBusiness string = "business"
	SubjectTypeEmployee string = "employee"

	EventLogin     string = "login"
	EventView      string = "view"
	EventCreate    string = "create"
	EventEdit      string = "edit"
	EventDelete    string = "delete"
	EventExport    string = "export"
	EventImport    string = "import"
	EventStatus    string = "status"
	EventStage     string = "stage"
	EventComment   string = "comment"
	EventManager   string = "manager"
	EventAnswer    string = "answer"
	EventMessenger string = "messenger"
	EventQueue     string = "queue"
	EventOther     string = "other"

	ObjectTypeCustomer    string = "customer"
	ObjectTypeRequest     string = "request"
	ObjectTypeStock       string = "stock"
	ObjectTypeSale        string = "sale"
	ObjectTypeTask        string = "task"
	ObjectTypeMessenger   string = "messenger"
	ObjectTypeRemind      string = "remind"
	ObjectTypeEmail       string = "email"
	ObjectTypeEmailsystem string = "emailsystem"
	ObjectTypeCall        string = "call"
	ObjectTypeSms         string = "sms"
	ObjectTypeDelivery    string = "delivery"
	ObjectTypeComment     string = "comment"
	ObjectTypeBlank       string = "blank"
	ObjectTypeApp         string = "app"
)

type whPayload interface {
	WHStockPayload
}

func WebhookHandler[T whPayload](payloadCh chan<- *T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		payload := new(T)
		if err := json.Unmarshal(body, payload); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		go func() {
			select {
			case payloadCh <- payload:
			case <-time.After(5 * time.Second):
			}
		}()

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
