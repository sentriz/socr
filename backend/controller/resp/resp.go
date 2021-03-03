package resp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.senan.xyz/socr/backend/db"
	"go.senan.xyz/socr/backend/hasher"
)

func Write(w http.ResponseWriter, body interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Response interface{} `json:"result"`
	}{
		Response: body,
	})
}

func Error(w http.ResponseWriter, status int, format string, a ...interface{}) {
	w.WriteHeader(status)
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		Error: fmt.Sprintf(format, a...),
	})
}

type ID struct {
	ID hasher.ID `json:"id"`
}

type Status struct {
	Status string `json:"status"`
}

type About struct {
	Version          string                          `json:"version"`
	APIKey           string                          `json:"api_key"`
	SocketClients    int                             `json:"socket_clients"`
	ImportPath       string                          `json:"import_path"`
	ScreenshotsPath  string                          `json:"screenshots_path"`
	ScreenshotsCount []db.CountDirectoriesByAliasRow `json:"screenshots_indexed"`
}

type Token struct {
	Token string `json:"token"`
}

type ImportStatus struct {
	Running bool `json:"running"`
}
