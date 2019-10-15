package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"kcapture/models"
	"kcapture/shared"

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
	//podMap := make(map[string][]models.PodInfo)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
		shared.LogMessage("ERROR", "error reading body")
	}

	json.Unmarshal(reqBody, &pods)
	json.NewEncoder(w).Encode(pods)

}

func main() {

	//Fetch the API PORT
	apiPort := shared.GetEnv("NODE_API_PORT", "9091")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", shared.HealthCheck)
	router.HandleFunc("/v1/nodeapi", nodeAPI).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+apiPort, router))
}
