package service

import (
	"fmt"
	"k8s_rbac/pkg/logger"
	"k8s_rbac/utilities"
	"net/http"
)

func ValidateEmailRequest(request utilities.RequestData) error {

	if len(request.Recipients) == 0 || request.Body == "" || request.Subject == "" {
		return fmt.Errorf("recipients, body and subject field is required")
	}
	return nil
}

func HandleError(err error, response *http.Response, w http.ResponseWriter) {
	if err != nil {
		logger.ErrorLogger.Printf("Failed to send message. Status Code: %v", response)
	} else {
		logger.InfoLogger.Printf("Message sent. Status Code: %v", response)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message sent"))
	}
}
