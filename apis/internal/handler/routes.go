package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"pod-transporter/apis/utils"
)

type RouteHandler struct {
	Config *rest.Config
}

type Input struct {
	SourceNamespace string
	DestNamespace   string
}

// HealthCheck use this for check it http util is working or not
func (rh *RouteHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("HealthCheck function started.")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (rh *RouteHandler) PodsReplicator(w http.ResponseWriter, r *http.Request) {
	cs := utils.GenerateClientSet(rh.Config)
	err := replicatePodsAcrossNamespaces(r, cs)
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
	}
	w.WriteHeader(http.StatusOK)
}

func replicatePodsAcrossNamespaces(r *http.Request, cs kubernetes.Interface) error {
	ctx := context.Background()
	var input Input

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Printf("error in decoding input")
		return err
	}
	err = utils.ReplicatePods(ctx, cs, input.SourceNamespace, input.DestNamespace, metav1.ListOptions{})
	if err != nil {
		log.Printf("error in replicating pods")
		return err
	}
	log.Printf("pods replicated successfully")
	return nil
}

func (rh *RouteHandler) PodLister(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cs := utils.GenerateClientSet(rh.Config)
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	podList, err := utils.ListPods(ctx, cs, namespace, metav1.ListOptions{})
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
	}
	jsonResp, err := json.Marshal(podList)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	w.WriteHeader(http.StatusOK)
}
