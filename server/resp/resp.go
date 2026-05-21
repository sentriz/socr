package resp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Write(w http.ResponseWriter, body any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Response any `json:"result"`
	}{
		Response: body,
	}); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func Errorf(w http.ResponseWriter, status int, format string, a ...any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: fmt.Sprintf(format, a...),
	}); err != nil {
		log.Printf("error encoding error response: %v", err)
	}
}
