package models

//Deployment struct for the request
type Deployment struct {
	Label     string `json:"label"`
	Namespace string `json:"namespace"`
}

//PodInfo to collect pod data
type PodInfo struct {
	//Deployment string `json:"deployment"`
	Name string `json:"name"`
	Node string `json:"node"`
	IP   string `json:"ip"`
}

//PodCollection to collect pod data
type PodCollection struct {
	// NodeName string
	Pods []PodInfo `json:"pods"`
}

//NodePodMap to store the pod to node mapping
type NodePodMap map[string][]PodInfo

//type NodePodMap map[string]*PodCollection
