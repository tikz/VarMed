package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// UniProtResponse represents some basic data from an UniProt accession.
// It is used to validate the user input in the New Job form.
type UniProtResponse struct {
	ID       string   `json:"id"`
	Sequence string   `json:"sequence"`
	PDBs     []string `json:"pdbs"`
	Name     string   `json:"name"`
	Gene     string   `json:"gene"`
	Organism string   `json:"organism"`
}

// UniProtEndpoint handles GET /api/uniprot/:unpID
func UniProtEndpoint(c *gin.Context) {
	id := c.Param("unpID")

	u, err := LoadUniProt(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UniProtResponse{
		ID:       u.ID,
		Sequence: u.Sequence,
		PDBs:     u.PDBIDs,
		Name:     u.Name,
		Gene:     u.Gene,
		Organism: u.Organism,
	})
}

func NewJobEndpoint(c *gin.Context) {
	id := c.Param("unpID")

	req := JobRequest{UniProtID: id}
	j := NewJob(&req)
	queue := c.MustGet("queue").(*Queue)

	// Add job to queue
	queue.Add(&j)
	fmt.Println("queue len", queue.Length())

	c.Data(http.StatusOK, "", []byte(j.ID))
}

// WSProcessEndpoint handles WebSocket /ws/:jobID
func WSProcessEndpoint(c *gin.Context) {
	id := c.Param("jobID")
	queue := c.MustGet("queue").(*Queue)
	job, _ := queue.GetJob(id)
	fmt.Println(job)
	wshandler(c.Writer, c.Request, job.MsgChan)
}

func wshandler(w http.ResponseWriter, r *http.Request, c <-chan string) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // TODO: remove

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade fail: %+v", err)
		return
	}

	for {
		select {
		case m := <-c:
			conn.WriteMessage(1, []byte(m))
		default:
		}
	}
}

func httpServe() {
	r := gin.Default()
	r.Use(cors.Default()) // TODO: remove

	queue := NewQueue(2)
	r.Use(func(c *gin.Context) {
		c.Set("queue", queue)
		c.Next()
	})

	r.Use(static.Serve("/", static.LocalFile("web/output", true)))
	// TODO: embed web/output inside binary

	// REST API endpoints
	r.GET("/api/uniprot/:unpID", UniProtEndpoint)
	r.GET("/api/new-job/:unpID", NewJobEndpoint)
	r.GET("/ws/:jobID", WSProcessEndpoint)

	log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	r.Run(":" + cfg.HTTPServer.Port)
}
