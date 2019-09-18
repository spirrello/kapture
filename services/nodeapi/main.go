package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/liaisontechnologies/kcapture/models"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("200 OK")
}

/*
nodeAPI receives the request and starts processing
*/
func nodeAPI(w http.ResponseWriter, r *http.Request) {
	var pods models.PodInfo
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
	}

	json.Unmarshal(reqBody, &pods)

	//json.NewEncoder(w).Encode(fetchPods(deploy.Label, deploy.Namespace))

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", healthCheck)
	router.HandleFunc("/v1/nodeapi", nodeAPI).Methods("POST")
	log.Fatal(http.ListenAndServe(":9091", router))
}
