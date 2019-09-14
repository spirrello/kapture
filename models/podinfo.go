package models

//PodInfo to collect pod data
type PodInfo struct {
	//Deployment string `json:"deployment"`
	Name string `json:"name"`
	Node string `json:"node"`
	IP   string `json:"ip"`
}
