package health

import (
	"net/http"
	"k8s_rbac/pkg/logger"
	"encoding/json"
)

type HealthResponse struct {
	Message string `json:"message"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := "Service is healthy"

	resp := HealthResponse{Message: msg}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		logger.ErrorLogger.Printf("Error happened in JSON marshal. Err: %s, Data: %+v", err, resp)
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		logger.ErrorLogger.Printf("Error writing response. Err: %s", err)
	}

	logger.InfoLogger.Printf("Service is healthy")
}