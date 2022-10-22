package main

import (
	"fmt"
	"log"
	"net/http"
	"pod-transporter/apis/internal/handler"
	"pod-transporter/apis/utils"

	"github.com/gorilla/mux"
)

func handleRequests() {
	// creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)
	// fetch cluster config
	cfg, err := utils.GetClusterConfig()
	if err != nil {
		log.Fatalf("failed to generate cluster config err: %v", err)
	}
	var routeHandler = handler.RouteHandler{
		Config: cfg,
	}
	replicatePath := "/api/v1/pods/replicate"
	listPath := "/api/v1/pods/namespace/{namespace}"
	router.HandleFunc(replicatePath, routeHandler.PodsReplicator)
	router.HandleFunc(listPath, routeHandler.PodLister)
	router.HandleFunc("/health", routeHandler.HealthCheck)
	router.HandleFunc("/", routeHandler.HealthCheck)
	err = http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatalf("error running the api server %v", err)
	}
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
