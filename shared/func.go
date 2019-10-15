package shared

import (
	"encoding/json"
	"kcapture/models"
	"log"
	"net/http"
	"os"
)

//GetEnv is a helper function for gathering env variables
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

//HealthCheck to run check.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(models.LogFormat{Loglevel: "info", Message: "200 OK"})
}

//LogMessage prints in JSON format
func LogMessage(logLevel, message string) {

	logStruct := models.LogFormat{Loglevel: logLevel, Message: message}
	logStr, _ := json.Marshal(logStruct)
	log.Println(string(logStr))

}
