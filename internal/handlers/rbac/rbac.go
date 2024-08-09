package rbacHandler

import (
	"k8s_rbac/pkg/certificate"
	"k8s_rbac/pkg/kuberclient"
	"k8s_rbac/pkg/logger"
	"encoding/json"

	"net/http"
)

func Rbac(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Username  string `json:"username"`
		Namespace string `json:"namespace"`
		Email     string `json:"email"`
	}
	
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	kuberclient.CreateRoles(requestData.Username, requestData.Namespace)

	caB64, certB64, keyB64, err := certificate.GenerateClientCreds(requestData.Username, "configs/ca.crt", "configs/ca.key")
	if err != nil {
		logger.ErrorLogger.Println("Error: Failed to generate client certificates ", err)
		http.Error(w, "Failed to generate client certificates", http.StatusInternalServerError)
		return
	}

	kuberclient.CreateKubeConfig(requestData.Username, caB64, certB64, keyB64)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("RBAC roles created"))

}

func ExtendRbac(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Username   string   `json:"username"`
		Namespaces []string `json:"namespaces"`
		Email      string   `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	for _, ns := range requestData.Namespaces {
		kuberclient.CreateRoles(requestData.Username, ns)

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("RBAC roles created, Access Extended"))
}
