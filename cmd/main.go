package main

import (
	"k8s_rbac/routes"

	"k8s_rbac/pkg/logger"
	"net/http"
)

func main() {
	logger.InitLogger()

	r := routes.SetUpRouter()

	http.ListenAndServe(":8080", (r))
}
