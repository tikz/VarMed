package main

import (
	"encoding/json"
	"log"
	"net/http"
	"varq/config"
	"varq/protein"
)

var (
	cfg *config.Config
)

// StatusResponse contains the JSON response about the server status
type StatusResponse struct {
	Code int    `json:"status_code"`
	Msg  string `json:"status_msg"`
}

// ErrorResponse contains the JSON response when an error occurs
type ErrorResponse struct {
	Msg string `json:"error_msg"`
}

func errorResponse(msg string) []byte {
	s := ErrorResponse{Msg: msg}
	out, _ := json.Marshal(s)
	return out
}

// proteinEndpoint is the function for the GET /proteinEndpoint endpoint
// Shows all parsed and calculated data for a given UniProt accession. For debug purposes only.
func proteinEndpoint(w http.ResponseWriter, r *http.Request) {
	params, ok := r.URL.Query()["uniprot"]

	if !ok || len(params[0]) < 1 {
		w.Write(errorResponse("GET params not present in request."))
		return
	}

	uniprotID := params[0]
	log.Println("New request from", r.RemoteAddr, "- UniProt", uniprotID)

	p, err := protein.NewProtein(uniprotID)
	if err != nil {
		w.Write(errorResponse(err.Error()))
		return
	}

	out, _ := json.Marshal(p)
	w.Write(out)
}

// statusEndpoint is the function for the GET /status endpoint
// Shows general information about the VarQ server status
func statusEndpoint(w http.ResponseWriter, r *http.Request) {
	s := StatusResponse{Code: 0, Msg: "online"}
	out, _ := json.Marshal(s)
	w.Write(out)
}

func init() {
	// Load config.yaml
	c, err := config.LoadFile("config.yaml")
	if err != nil {
		log.Fatalf("Cannot open and parse config.yaml: %v", err)
	}
	cfg = c
}

func main() {
	// !DEBUG
	// _, err := protein.NewProtein("P69892")
	// fmt.Println("newp err:", err)

	// REST API entrypoints
	http.HandleFunc("/status", statusEndpoint)
	http.HandleFunc("/protein", proteinEndpoint)

	log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	http.ListenAndServe(":"+cfg.HTTPServer.Port, nil)
}
