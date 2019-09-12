package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//Pods struct to collect deployment pods
type Pods struct {
	Deployment string `json:"deployment"`
	Namespace  string `json:"namespace"`
}

//Deployment struct for the request
type Deployment struct {
	Label     string `json:"label"`
	Namespace string `json:"namespace"`
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

//externalKubeClient creates the external cluster config
func externalKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil

}

//internalKubeClient creates the in-cluster config
func internalKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func setupKubeClient() (*kubernetes.Clientset, error) {
	if os.Getenv("KUBECONFIG") != "" {
		clientset, err := externalKubeClient(os.Getenv("KUBECONFIG"))
		if err != nil {
			return nil, err
		}
		return clientset, nil
	}

	clientset, err := internalKubeClient()
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func fetchPods(label string, namespace string) {

	//setup connection to kube API
	clientset, err := setupKubeClient()
	if err != nil {
		panic(err.Error())
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: label})
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Spec.NodeName)
	}

	if errors.IsNotFound(err) {
		fmt.Printf("Pod not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}

}

func deployment(w http.ResponseWriter, r *http.Request) {
	var deployment Deployment
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
	}

	json.Unmarshal(reqBody, &deployment)

	json.NewEncoder(w).Encode(deployment.Label)

	fetchPods(deployment.Label, deployment.Namespace)

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", healthCheck)
	router.HandleFunc("/v1/deployment", deployment).Methods("POST")
	log.Fatal(http.ListenAndServe(":9090", router))
}
