package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// StatusEndpoint handles GET /api/status
// Returns the API status.
func StatusEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "online"})
}

// UniProtEndpoint handles GET /api/uniprot/:unpID
// Fetches and returns fields from an UniProt entry.
func UniProtEndpoint(c *gin.Context) {
	id := c.Param("unpID")

	u, err := loadUniProt(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, u)
}

// JobEndpoint handles GET /api/job/:jobID
// Returns general status and request info about a job.
func JobEndpoint(c *gin.Context) {
	id := c.Param("jobID")
	queue := c.MustGet("queue").(*Queue)

	// From queue
	job, err := queue.GetJob(id)
	if err == nil {
		c.JSON(http.StatusOK, job)
		return
	}

	// From file
	job, err = loadJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

// JobPDBEndpoint handles GET /api/job/:jobID/:pdbID
// Returns results about a structure in a job.
func JobPDBEndpoint(c *gin.Context) {
	jobID := c.Param("jobID")
	pdbID := c.Param("pdbID")

	job, err := loadJob(jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if _, ok := job.Pipeline.Results[pdbID]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, job.Pipeline.Results[pdbID])
}

// NewJobEndpoint handles POST /api/new-job
// Starts a new job.
func NewJobEndpoint(c *gin.Context) {
	req := JobRequest{}
	c.BindJSON(&req)
	req.IP = c.ClientIP()
	req.Time = time.Now()
	req.SAS = []string{"M1K"}
	// TODO: data should be already client side validated, but
	// do it again server side

	j := NewJob(&req)
	queue := c.MustGet("queue").(*Queue)
	queue.Add(&j)

	c.JSON(http.StatusOK, gin.H{"id": j.ID, "error": ""})
}

// WSProcessEndpoint handles WebSocket /ws/:jobID
func WSProcessEndpoint(c *gin.Context) {
	id := c.Param("jobID")
	queue := c.MustGet("queue").(*Queue)
	job, err := queue.GetJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	wsHandler(c.Writer, c.Request, job)
}

func wsHandler(w http.ResponseWriter, r *http.Request, j *Job) {
	upg := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	upg.CheckOrigin = func(r *http.Request) bool { return true } // TODO: remove in production, unsafe

	ws, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade fail: %+v", err)
		return
	}

	msgTicker := time.NewTicker(100 * time.Millisecond)
	defer func() {
		msgTicker.Stop()
		ws.Close()
	}()

	i := len(j.msgs)
	for {
		select {
		case <-msgTicker.C:
			if i < len(j.msgs) {
				msg := j.msgs[i]
				ws.WriteMessage(1, []byte(msg))
				i++

				if j.Status == statusDone ||
					j.Status == statusSaved ||
					j.Status == statusError {
					return
				}
			}
		}

	}
}

// CIFEndpoint handles GET /api/structure/cif/:pdbID
func CIFEndpoint(c *gin.Context) {
	id := c.Param("pdbID")

	p, err := readPDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/plain", p.RawCIF)
}

func httpServe() {
	r := gin.Default()
	r.Use(cors.Default()) // TODO: remove in production, unsafe

	// Job queue, pass inside context to Gin methods
	queue := NewQueue(cfg.VarQ.JobWorkers)
	r.Use(func(c *gin.Context) {
		c.Set("queue", queue)
		c.Next()
	})

	r.Use(static.Serve("/", static.LocalFile("web/output", true)))
	// TODO: embed web/output files inside binary

	// API endpoints
	r.GET("/api/status", StatusEndpoint)
	r.GET("/api/uniprot/:unpID", UniProtEndpoint)
	r.GET("/api/job/:jobID", JobEndpoint)
	r.GET("/api/job/:jobID/:pdbID", JobPDBEndpoint)
	r.GET("/api/structure/cif/:pdbID", CIFEndpoint)
	r.GET("/ws/:jobID", WSProcessEndpoint)

	r.POST("/api/new-job", NewJobEndpoint)

	// Let React Router manage all root paths not declared here
	r.NoRoute(func(c *gin.Context) {
		c.File("web/output/index.html")
	})

	log.Printf("Starting VarQ web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	r.Run(":" + cfg.HTTPServer.Port)
}
