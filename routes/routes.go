package routes

import (
	healthHandler "k8s_rbac/internal/handlers/health"
	rbacHandler "k8s_rbac/internal/handlers/rbac"
	"k8s_rbac/pkg/logger"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	logger.InitLogger()
	r := mux.NewRouter()

	r.HandleFunc("/rbac", rbacHandler.Rbac).Methods("POST")
	r.HandleFunc("/extend_rbac", rbacHandler.ExtendRbac).Methods("POST")
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")

	logger.InfoLogger.Println("Server listening on 8080..")
	return r
}

