package main

import (
	"encoding/json"
	"net/http"
)

type Status struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func main() {
	// API entrypoints
	http.HandleFunc("/status", status)

	http.ListenAndServe(":3000", nil)
}

// Status is the entrypoint function for GET /status
// Displays general information about the VarQ server
func status(w http.ResponseWriter, r *http.Request) {
	s := Status{StatusCode: 0, StatusMsg: "Online"}
	out, _ := json.Marshal(s)
	w.Write(out)
}
