package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/liaisontechnologies/kcapture/models"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//podStruct to collect pod data
// type podStruct struct {
// 	//Deployment string `json:"deployment"`
// 	Name string `json:"name"`
// 	Node string `json:"node"`
// 	IP   string `json:"ip"`
// }

//deployment struct for the request
type deployment struct {
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

func fetchPods(label string, namespace string) []models.PodInfo {

	//setup connection to kube API
	clientset, err := setupKubeClient()
	if err != nil {
		panic(err.Error())
	}

	//Setup podSlice
	podSlice := []models.PodInfo{}

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: label})
	if errors.IsNotFound(err) {
		log.Fatal("Pod not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Fatalf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}
	// capturePods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "k8s-app=filebeat"})
	for _, pod := range pods.Items {
		podSlice = append(podSlice, models.PodInfo{Name: pod.Name, Node: pod.Spec.NodeName, IP: pod.Status.PodIP})
		//LETS fetch the capture pod per node
		//k8-app = filebeat
		//os.Getenv("KCAPTURE_CAPTURE")
		capturePod, err := clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{
			LabelSelector: "k8s-app=filebeat", FieldSelector: "spec.nodeName=" + pod.Spec.NodeName,
		})
		if err != nil {
			panic(err.Error())
		}
		log.Println(capturePod.Items)

		// podInfo := models.PodInfo{Name: pod.Name, Node: pod.Spec.NodeName, IP: pod.Status.PodIP}
		// podCollection.NodeName
		// podCollection.NodeName.Pods = append(podCollection.NodeName.Pods, podInfo)
	}

	log.Println(podSlice)
	return podSlice

}

// func fetchNodes(podCollection models.PodCollection) {

// 	//setup connection to kube API
// 	clientset, err := setupKubeClient()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	//var podInfo models.PodInfo

// 	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})

// 	for _, node := range nodes.Items {
// 		log.Println(node.Name)
// 		podCollection = append(podCollection, NodeName: node.Name,)
// 	}
// }

/*
pods receives the request and calls fetchPods to
retrieve pod details
*/
func pods(w http.ResponseWriter, r *http.Request) {
	//initialize deploy and podCollection structs first
	var deploy deployment
	// var podCollection models.PodCollection
	// fetchNodes(*podCollection)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
	}

	json.Unmarshal(reqBody, &deploy)

	json.NewEncoder(w).Encode(fetchPods(deploy.Label, deploy.Namespace))

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", healthCheck)
	router.HandleFunc("/v1/pods", pods).Methods("POST")
	log.Fatal(http.ListenAndServe(":9090", router))
}
