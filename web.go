package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

// UniProtEndpoint is the function for the /uniprot/{id} endpoint
// Shows parsed and calculated data for a given UniProt accession for debug purposes.
func UniProtEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]
	log.Println("New request from", r.RemoteAddr, "- UniProt", ID)

	p, err := RunPipeline(ID, []string{})
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

func httpServe() {
	// REST API entrypoints
	r := mux.NewRouter()
	r.HandleFunc("/status", statusEndpoint)
	r.HandleFunc("/uniprot/{ID}", UniProtEndpoint)
	http.Handle("/", r)

	log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	http.ListenAndServe(":"+cfg.HTTPServer.Port, nil)
}
