package main

import (
	"encoding/json"
	"net/http"
)

// JSON response convenience function
func JSONResponse(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Add("Content-Type", "application/json")

	// Error Response - Return early
	if err != nil {
		w.WriteHeader(status)
		return
	}

	// Try to handle data
	jRes, err := json.Marshal(data)
	if err != nil {
		jErr, _ := json.Marshal(map[string]interface{}{
			"Data Error": err.Error(),
		})
		w.WriteHeader(500)
		w.Write(jErr)
	}
	w.WriteHeader(status)
	w.Write(jRes)
}
