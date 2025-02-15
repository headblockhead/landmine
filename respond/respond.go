package respond

import (
	"encoding/json"
	"net/http"
)

func WithJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		json.NewEncoder(w).Encode(struct {
			Message    string `json:"message"`
			StatusCode int    `json:"status_code"`
		}{
			Message:    "failed to encode json",
			StatusCode: http.StatusInternalServerError,
		})
	}
}
