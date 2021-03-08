package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Write(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Response interface{} `json:"result"`
	}{
		Response: body,
	})
}

func Error(w http.ResponseWriter, status int, format string, a ...interface{}) {
	w.WriteHeader(status)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: fmt.Sprintf(format, a...),
	})
}
