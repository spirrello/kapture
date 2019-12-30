package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"kapture/models"
	"kapture/shared"

	"github.com/gorilla/mux"
)

func main() {

	//Fetch the API PORT
	apiPort := shared.GetEnv("NODE_API_PORT", "9091")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", shared.HealthCheck)
	router.HandleFunc("/v1/nodeapi", nodeAPI).Methods("POST")
	http.ListenAndServe(":"+apiPort, router)
}

/*
nodeAPI receives the request and starts processing
*/
func nodeAPI(w http.ResponseWriter, r *http.Request) {
	var pods []models.PodInfo
	//podMap := make(map[string][]models.PodInfo)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
		shared.LogMessage("ERROR", "error reading body")
	}

	json.Unmarshal(reqBody, &pods)
	//json.NewEncoder(w).Encode(pods)

	ctx := context.Background()
	startCapture(ctx)

}

//runCapture invokes packetCapture with a channel
func startCapture(ctx context.Context) {

	shared.LogMessage("INFO", "new capture starting")
	// Create a channel for signal handling
	c := make(chan struct{})
	defaultCapTimeout := shared.GetEnv("NODE_API_CAP_TIMEOUT", "10")
	// Define a cancellation after 1s in the context
	capTimeout, _ := strconv.Atoi(defaultCapTimeout)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(capTimeout+5)*time.Second)
	defer cancel()

	go func() {
		packetCapture(c)
	}()

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	case <-c:
		shared.LogMessage("INFO", "capture completed successfully")

	}

}

//packetCapture executes the packet capture
func packetCapture(c chan struct{}) {

	//default ENV values
	defaultNIC := shared.GetEnv("NODE_API_NIC", "en0")
	defaultFullPacket := shared.GetEnv("NODE_API_FULL_PACKET", "0")
	defaultCapTimeout := shared.GetEnv("NODE_API_CAP_TIMEOUT", "10")
	defaultTotalCaptures := shared.GetEnv("NODE_API_TOTAL_CAPTURES", "1")
	defaultCaptureFile := shared.GetEnv("NODE_API_CAPTURE_FILE", "test.pcap")

	shared.LogMessage("INFO", "capture started")
	cmd := exec.Command("tcpdump", "-i", defaultNIC, "-nn", "-s", defaultFullPacket, "-G", defaultCapTimeout, "-W", defaultTotalCaptures, "-w", defaultCaptureFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		shared.LogMessage("ERROR", string("cmd.Run() failed with"+err.Error()))
	}
	shared.LogMessage("INFO", string(out))

	c <- struct{}{}

}
