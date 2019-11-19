package models

//LogFormat struct for return log messages in json format
type LogFormat struct {
	Loglevel string `json:"level"`
	Caller   string `json:"caller"`
	Message  string `json:"message"`
}

//Deployment struct for the initial request
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

//CaptureInstructions organizes the instructions
//for the node API
type CaptureInstructions struct {
	State string `json:"state"`
}
