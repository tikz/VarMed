package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/tikz/bio"
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

	u, err := bio.LoadUniProt(id)
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

// JobPDBCSVEndpoint handles GET /api/job/:jobID/:pdbID/csv
// Returns a CSV file of variants in a job.
func JobPDBCSVEndpoint(c *gin.Context) {
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
	results := job.Pipeline.Results[pdbID]

	filename := fmt.Sprintf("%s_%s_%s.csv", results.UniProt.ID, results.PDB.ID, jobID[:5])
	c.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.String(http.StatusOK, ResultsCSV(job.Pipeline.Results[pdbID]))
}

// NewJobEndpoint handles POST /api/new-job
// Starts a new job.
func NewJobEndpoint(c *gin.Context) {
	req := JobRequest{}
	c.BindJSON(&req)
	req.IP = c.ClientIP()
	req.Time = time.Now()

	// Check if job already exists
	j, err := loadJob(generateID(&req))
	if err != nil {
		j = NewJob(&req)
		queue := c.MustGet("queue").(*Queue)
		queue.Add(j)
	}

	c.JSON(http.StatusOK, gin.H{"id": j.ID, "error": ""})
}

// CIFEndpoint handles GET /api/structure/cif/:pdbID
func CIFEndpoint(c *gin.Context) {
	id := c.Param("pdbID")

	p, err := bio.LoadPDB(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cif, err := p.RawCIF()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/plain", cif)
}

// MutatedPDBEndpoint handles GET /api/mutated/:pdbID/:mutation
func MutatedPDBEndpoint(c *gin.Context) {
	pdbID := c.Param("pdbID")
	mutation := c.Param("mutation")

	pdbPath := fmt.Sprintf("%s/%s/%s/%s_Repair_1.pdb",
		cfg.Paths.FoldXMutations, pdbID, mutation, pdbID)

	if _, err := os.Stat(pdbPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	pdb, err := ioutil.ReadFile(pdbPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/plain", pdb)
}

func httpServe() {
	r := gin.Default()
	r.Use(cors.Default()) // TODO: remove in production, unsafe

	// Job queue, pass inside context to Gin methods
	queue := NewQueue(cfg.VarMed.JobWorkers)
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
	r.GET("/api/job/:jobID/:pdbID/csv", JobPDBCSVEndpoint)
	r.GET("/api/structure/cif/:pdbID", CIFEndpoint)
	r.GET("/api/mutated/:pdbID/:mutation", MutatedPDBEndpoint)
	r.GET("/ws/job/:jobID", WSJobEndpoint)
	r.GET("/ws/queue", WSQueueEndpoint)

	r.POST("/api/new-job", NewJobEndpoint)

	// Let React Router manage all root paths not declared here
	r.NoRoute(func(c *gin.Context) {
		c.File("web/output/index.html")
	})

	log.Printf("Starting VarMed web server: http://127.0.0.1:%s/", cfg.HTTPServer.Port)
	r.Run(":" + cfg.HTTPServer.Port)
}
