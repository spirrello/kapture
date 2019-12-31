package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"

	"kapture/models"
	"kapture/shared"

	"github.com/gorilla/mux"
)

// //healthCheck to run check.
// func healthCheck(w http.ResponseWriter, r *http.Request) {
// 	json.NewEncoder(w).Encode(models.LogFormat{Loglevel: "info", Message: "200 OK"})
// }

//nodeAPI receives the request and starts processing
func nodeAPI(w http.ResponseWriter, r *http.Request) {
	var pods []models.PodInfo
	//podMap := make(map[string][]models.PodInfo)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
		shared.LogMessage("ERROR", "error reading body")
	}

	json.Unmarshal(reqBody, &pods)
	json.NewEncoder(w).Encode(pods)

	for _, pod := range pods {
		//log pod info that will be captured and invoke goroutine
		podString, _ := json.Marshal(pod)
		shared.LogMessage("INFO", string(podString))
		go runCapture(pod)
	}

}

//runCapture executes packet captures
func runCapture(pod models.PodInfo) {
	//default ENV values
	defaultNIC := shared.GetEnv("NODE_API_NIC", "any")
	defaultFullPacket := shared.GetEnv("NODE_API_FULL_PACKET", "0")
	var capTimeout string
	if pod.Time == "" {
		capTimeout = shared.GetEnv("NODE_API_CAP_TIMEOUT", "60")
	} else {
		capTimeout = pod.Time
	}
	var totalCaptures string
	if pod.TotalCaptures == "" {
		totalCaptures = shared.GetEnv("NODE_API_TOTAL_CAPTURES", "1")
	} else {
		totalCaptures = pod.TotalCaptures
	}
	defaultPcapDir := shared.GetEnv("NODE_API_PCAP_DIR", "pcap")
	pcapFile := defaultPcapDir + "/" + pod.Name + ".pcap"
	//run tcpdump
	cmd := exec.Command("tcpdump", "-i", defaultNIC, "-nn", "-s", defaultFullPacket, "-G", capTimeout, "-W", totalCaptures, "-w", pcapFile, "host", pod.IP)
	out, err := cmd.CombinedOutput()
	if err != nil {
		shared.LogMessage("ERROR", string("cmd.Run() failed with"+err.Error()))
		shared.LogMessage("ERROR", "defaultNIC: "+defaultNIC+" defaultFullPacket: "+defaultFullPacket+"capTimeout: "+capTimeout+" totalCaptures:"+totalCaptures)
	}
	shared.LogMessage("INFO", string(out))
}

func main() {

	//Fetch the API PORT
	apiPort := shared.GetEnv("NODE_API_PORT", "9091")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", shared.HealthCheck)
	router.HandleFunc("/v1/nodeapi", nodeAPI).Methods("POST")
	http.ListenAndServe(":"+apiPort, router)
}
