package utils

import (
	"net/http"
	"encoding/json"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err error, message string) {
    errorResp := ErrorResponse{
        Error:   err.Error(),
        Message: message,
    }
    WriteJSON(w, status, errorResp)
}

func ParseJSON(r *http.Request, dest interface{}) error {
    return json.NewDecoder(r.Body).Decode(dest)
}