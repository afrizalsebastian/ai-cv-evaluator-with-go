package api

import (
	"encoding/json"
	"net/http"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, response models.WebResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
