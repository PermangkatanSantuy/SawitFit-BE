package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ParseJSON: Membaca data JSON yang dikirim dari Frontend
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

// WriteJSON: Mengirim balasan (Response) ke Frontend dalam format JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// WriteError: Mengirim pesan error yang rapi ke Frontend
func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}