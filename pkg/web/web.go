// pkg/web/web.go
package web

import (
	"encoding/json"
	"net/http"
)

// Success envía una respuesta HTTP con código de estado y payload JSON.
func Success(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// Error envía una respuesta HTTP con código de error y mensaje JSON.
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := map[string]string{
		"error": message,
	}

	json.NewEncoder(w).Encode(resp)
}
