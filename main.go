package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//"github.com/kubernetes/client-go/rest"
)

//Pods struct to collect deployment pods
type Pods struct {
	Deployment string `json:"deployment"`
	Namespace  string `json:"namespace"`
}

//LogFormat struct for return log messages in json format
type LogFormat struct {
	Loglevel string `json:"level"`
	Message  string `json:"message"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("200 OK")
}

func logMessage(level string, message string) {
	var logContent LogFormat

	json.Unmarshal([]byte(message), &logContent)

	log.Println(logContent)
}

//k8sClient sets up the K8s api client
func k8sClient() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func pods(w http.ResponseWriter, r *http.Request) {
	var pod Pods
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
	}

	json.Unmarshal(reqBody, &pod)

	json.NewEncoder(w).Encode(pod)

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", healthCheck)
	router.HandleFunc("/v1/pods", pods).Methods("POST")
	log.Fatal(http.ListenAndServe(":9090", router))
}
