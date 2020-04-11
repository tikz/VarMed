package main

import (
	"encoding/json"
	"log"
	"net/http"
	"varq/uniprot"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type StatusResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
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

type UniProtResponse struct {
	Sequence string   `json:"sequence"`
	PDBs     []string `json:"pdbs"`
	Name     string   `json:"name"`
	Gene     string   `json:"gene"`
	Organism string   `json:"organism"`
}

// UniProtEndpoint handles GET /api/uniprot/:id
func UniProtEndpoint(c *gin.Context) {
	id := c.Param("id")

	u, err := uniprot.NewUniProt(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UniProtResponse{
		Sequence: u.Sequence,
		PDBs:     u.PDBIDs,
		Name:     u.Name,
		Gene:     u.Gene,
		Organism: u.Organism,
	})
}

// statusEndpoint is the function for the GET /status endpoint
// Shows general information about the VarQ server status
// func statusEndpoint(w http.ResponseWriter, r *http.Request) {
// 	s := StatusResponse{Code: 0, Msg: "online"}
// 	out, _ := json.Marshal(s)
// 	w.Write(out)
// }

func httpServe() {
	r := gin.Default()
	r.Use(cors.Default()) // TODO: remove

	// REST API entrypoints
	// r.HandleFunc("/status", statusEndpoint)
	r.GET("/api/uniprot/:id", UniProtEndpoint)

	log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	r.Run(":" + cfg.HTTPServer.Port)
	// r.PathPrefix("/output/").Handler(http.FileServer(http.Dir("./web/output/")))
	// http.Handle("/", r)

	// fs := http.FileServer(http.Dir("./web/output"))
	// http.Handle("/", fs)

	// log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	// err := http.ListenAndServe(":"+cfg.HTTPServer.Port, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
