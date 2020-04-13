package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"strconv"
	"time"
)

const (
	statusPending = 0
	statusProcess = 1
	statusDone    = 2
)

type JobRequest struct {
	IP        string
	Email     string
	Time      time.Time
	UniProtID string
	PDBIDs    []string
	ClinVar   bool
	Variants  map[int64]string
}

type Job struct {
	ID       string
	Request  *JobRequest
	Status   int
	MsgChan  chan string
	Started  time.Time
	Ended    time.Time
	Pipeline *Pipeline
	Error    error
}

// generateID returns a SHA256 hash of timestamp+UniProtID+256 random bits.
func (j *Job) generateID() string {
	rb := [32]byte{}
	_, err := rand.Read(rb[:])
	if err != nil {
		log.Fatal(err)
	}

	timestamp := []byte(strconv.FormatInt(j.Request.Time.Unix(), 10))
	unpID := []byte(j.Request.UniProtID)
	b := bytes.Join([][]byte{timestamp, unpID, rb[:]}, []byte("-"))
	hash := sha256.Sum256(b)

	return hex.EncodeToString(hash[:])
}

func NewJob(request *JobRequest) Job {
	j := Job{Request: request}
	j.MsgChan = make(chan string, 100)

	j.ID = j.generateID()

	return j
}

func (j *Job) Process() {
	j.Status = statusProcess
	j.Started = time.Now()

	j.Pipeline, _ = NewPipeline(j.Request.UniProtID, []string{"3CON", "5UHV", "6E6H"}, j.MsgChan)
	_ = j.Pipeline.RunPipeline()
	// j.PDBAnalyses, _ = RunPipeline("P01111", []string{"3CON"})
	j.Ended = time.Now()
}
