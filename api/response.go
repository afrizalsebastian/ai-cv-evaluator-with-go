package api

import (
	"encoding/json"
	"net/http"
)

type WebResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func CreateWebResponse(message string, status int, data interface{}, errors interface{}) WebResponse {
	return WebResponse{
		Message: message,
		Status:  status,
		Data:    data,
		Errors:  errors,
	}
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, response WebResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
