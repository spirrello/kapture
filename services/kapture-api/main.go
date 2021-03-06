package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"kapture/models"
	"kapture/shared"

	"github.com/gorilla/mux"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// //healthCheck to run check.
// func healthCheck(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(models.LogFormat{Loglevel: "info", Message: "200 OK"})
// }

// //logMessage prints in JSON format
// func logMessage(logLevel, message string) {

// 	logStruct := models.LogFormat{Loglevel: logLevel, Message: message}
// 	logStr, _ := json.Marshal(logStruct)
// 	log.Println(string(logStr))

// }

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

/*
fetchPods queries kubeapi to find the requested pods and returns a map of
them with the capture pods.
*/
func fetchPods(label string, namespace string) map[string][]models.PodInfo {

	//setup connection to kube API
	clientset, err := setupKubeClient()
	if err != nil {
		panic(err.Error())
	}

	//fetch pods
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: label})
	if len(pods.Items) == 0 {
		shared.LogMessage("ERROR", label+" not found")
	}
	if errors.IsNotFound(err) {
		log.Fatal("Pod not found\n")
		//logMessage("ERROR", "Pod not found")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		log.Fatalf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	}

	//This map uses the capture pods as the keys
	podMap := make(map[string][]models.PodInfo)
	for _, pod := range pods.Items {
		//LETS fetch the capture pod that exist on the same node as the requested pod.
		capturePods, err := clientset.CoreV1().Pods(os.Getenv("CAPTURE_PODS_NAMESPACE")).List(metav1.ListOptions{
			LabelSelector: os.Getenv("CAPTURE_PODS"), FieldSelector: "spec.nodeName=" + pod.Spec.NodeName,
		})
		if err != nil {
			panic(err.Error())
		} else {
			log.Println(label + " deployment found")

		}
		//podMap[capturePods.Items[0].Name] = append(podMap[capturePods.Items[0].Name], models.PodInfo{Name: pod.Name, IP: pod.Status.PodIP})
		podMap[capturePods.Items[0].Status.PodIP] = append(podMap[capturePods.Items[0].Status.PodIP], models.PodInfo{Name: pod.Name, IP: pod.Status.PodIP})
	}
	return podMap
}

//nodeInstruct posts to the nodeAPI with instructions.
func nodeInstruct(podMap map[string][]models.PodInfo) (string, error) {
	nodeAPIPort := shared.GetEnv("NODE_API_PORT", "9090")
	for k, v := range podMap {
		bytesRepresentation, err := json.Marshal(v)
		_, err = http.Post("http://"+k+":"+nodeAPIPort+"/v1/nodeapi", "application/json", bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			shared.LogMessage("ERROR", err.Error())
			return "", err
		}
	}
	return "200 OK", nil
}

/*
pods receives the request and calls fetchPods to
retrieve pod details
*/
func pods(w http.ResponseWriter, r *http.Request) {
	//initialize deploy and podCollection structs first
	var deploy models.Deployment

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
	}

	json.Unmarshal(reqBody, &deploy)

	podMap := fetchPods(deploy.Label, deploy.Namespace)
	json.NewEncoder(w).Encode(podMap)

	nodeInstructResult, err := nodeInstruct(podMap)
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(nodeInstructResult)
	}
}

func main() {

	//Fetch the API PORT
	//apiPort := shared.GetEnv("KAPTURE_API_PORT", "9090")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", shared.HealthCheck)
	router.HandleFunc("/v1/pods", pods).Methods("POST")
	log.Fatal(http.ListenAndServe(":9090", router))
}
