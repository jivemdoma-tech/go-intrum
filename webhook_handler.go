package gointrum

import (
	"encoding/json"
	"io"
	"net/http"
)

func WebhookHandler(callback func(WebhookEvent)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var event WebhookEvent
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		callback(event)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
