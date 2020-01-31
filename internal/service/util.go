package service

import (
	"encoding/json"
	"net/http"

	"github.com/lufanbb/Golang_REST_API_JWT_Auth/internal/model"
)

// RespondWithError takes the status code and error and write back the error json
func RespondWithError(w http.ResponseWriter, statusCode int, error model.Error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(error)
}

// RespondWithJSON takes the status code and error and write back the error json
func RespondWithJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
