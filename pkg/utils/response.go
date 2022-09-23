package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// ResponseWithJson writes a json response.
func ResponseWithJson(w http.ResponseWriter, status int, object any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		log.Err(err).Msg("Failed to encode json")
		return
	}
}
