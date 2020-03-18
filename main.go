package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// Config holds the parsed config file
var Config *Configuration

// StatusResponse contains the JSON response about the server status
type StatusResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// loadConfig opens and parses config.yaml
func loadConfig() {
	f, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error opening config.yaml: %v", err)
	}

	cfg := Configuration{}
	err = yaml.Unmarshal([]byte(f), &cfg)
	if err != nil {
		log.Fatalf("error parsing config.yaml: %v", err)
	}
	Config = &cfg
}

// status is the function for the GET /status entrypoint
// Shows general information about the VarQ server status
func status(w http.ResponseWriter, r *http.Request) {
	s := StatusResponse{StatusCode: 0, StatusMsg: "online"}
	out, _ := json.Marshal(s)
	w.Write(out)
}

func main() {
	loadConfig()

	// REST API entrypoints
	http.HandleFunc("/status", status)

	log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", Config.HTTPServer.Port)
	http.ListenAndServe(":"+Config.HTTPServer.Port, nil)
}
