package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"kcapture/models"
)

//healthCheck to run check.
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.LogFormat{Loglevel: "info", Message: "200 OK"})
}

/*
nodeAPI receives the request and starts processing
*/
func nodeAPI(w http.ResponseWriter, r *http.Request) {
	//var pods []models.PodInfo
	podMap := make(map[string][]models.PodInfo)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		json.NewEncoder(w).Encode("error reading body")
	}

	json.Unmarshal(reqBody, &podMap)
	json.NewEncoder(w).Encode(podMap)

	log.Println(podMap)

}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/healthcheck", healthCheck)
	router.HandleFunc("/v1/nodeapi", nodeAPI).Methods("POST")
	log.Fatal(http.ListenAndServe(":9091", router))
}
