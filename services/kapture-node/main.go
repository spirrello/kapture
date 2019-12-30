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

/*
nodeAPI receives the request and starts processing
*/
func nodeAPI(w http.ResponseWriter, r *http.Request) {
	var pods []models.PodInfo
	podMap := make(map[string][]models.PodInfo)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
		shared.LogMessage("ERROR", "error reading body")
	}

	json.Unmarshal(reqBody, &pods)
	json.NewEncoder(w).Encode(pods)

	go runCapture()

}

//runCapture executes packet captures
func runCapture() {
	//default ENV values
	defaultNIC := shared.GetEnv("NODE_API_NIC", "en0")
	defaultFullPacket := shared.GetEnv("NODE_API_FULL_PACKET", "0")
	defaultCapTimeout := shared.GetEnv("NODE_API_CAP_TIMEOUT", "10")
	defaultTotalCaptures := shared.GetEnv("NODE_API_TOTAL_CAPTURES", "1")
	defaultCaptureFile := shared.GetEnv("NODE_API_CAPTURE_FILE", "test.pcap")

	cmd := exec.Command("tcpdump", "-i", defaultNIC, "-nn", "-s", defaultFullPacket, "-G", defaultCapTimeout, "-W", defaultTotalCaptures, "-w", defaultCaptureFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		shared.LogMessage("ERROR", string("cmd.Run() failed with"+err.Error()))
	}
	shared.LogMessage("INFO", string(out))

	// select {
	// case <-ctx.Done():
	// 	fmt.Println(ctx.Err())
	// case <-c:
	// 	fmt.Println("Unexpected success!")
	// }
	// fmt.Println("capture completed successfully")

}

func main() {

	//Fetch the API PORT
	apiPort := shared.GetEnv("NODE_API_PORT", "9091")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", shared.HealthCheck)
	router.HandleFunc("/v1/nodeapi", nodeAPI).Methods("POST")
	http.ListenAndServe(":"+apiPort, router)
}
