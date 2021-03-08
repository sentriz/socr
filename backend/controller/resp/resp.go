package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Write(w http.ResponseWriter, body interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Response interface{} `json:"result"`
	}{
		Response: body,
	})
}

func Error(w http.ResponseWriter, status int, format string, a ...interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: fmt.Sprintf(format, a...),
	})
}
