package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

const (
	statusPending = 0
	statusProcess = 1
	statusDone    = 2
	statusSaved   = 3
	statusError   = 4
)

// JobRequest represents a job request from an user.
// Contains the user input and additional details.
type JobRequest struct {
	Name          string    `json:"name"`
	UniProtID     string    `json:"uniprotId"`
	PDBIDs        []string  `json:"pdbIds"`
	ClinVar       bool      `json:"clinvar"`
	VariationsPos []int     `json:"variationsPos"`
	VariationsAA  []string  `json:"variationsAa"`
	IP            string    `json:"ip"`
	Email         string    `json:"email"`
	Time          time.Time `json:"time"`
}

// Job represents the input and outputs of a single job ran by the pipeline.
type Job struct {
	ID       string      `json:"id"`
	Request  *JobRequest `json:"request"`
	Pipeline *Pipeline   `json:"-"`
	Status   int         `json:"status"`
	Started  time.Time   `json:"started"`
	Ended    time.Time   `json:"ended"`

	msgs  []string
	Error error `json:"-"`
}

// generateID returns a SHA256 hash of UniProtID+joined PDBIDs
func (j *Job) generateID() string { // TODO: include variations in hash after implementing that

	unpID := []byte(j.Request.UniProtID)
	pdbIDs := []byte(strings.Join(j.Request.PDBIDs, ""))
	b := bytes.Join([][]byte{unpID, pdbIDs}, []byte(""))
	hash := sha256.Sum256(b)

	return hex.EncodeToString(hash[:])
}

// NewJob returns a new job instance.
func NewJob(request *JobRequest) Job {
	j := Job{Request: request}

	j.ID = j.generateID()

	return j
}

// Process runs the pipeline for the job.
func (j *Job) Process() {
	j.Status = statusProcess
	j.Started = time.Now()

	msgChan := make(chan string, 100)
	j.Pipeline, _ = NewPipeline(j.Request.UniProtID, j.Request.PDBIDs, msgChan)

	go func() {
		for m := range msgChan {
			j.msgs = append(j.msgs, m)
		}
	}()

	err := j.Pipeline.Run()
	if err != nil {
		j.fail(err)
		return
	}

	j.Ended = time.Now()
	j.Status = statusDone

	err = writeJob(j)
	if err != nil {
		panic(err)
	}
	j.Status = statusSaved
}

// fail handles the given error message and updates the status.
func (j *Job) fail(err error) {
	j.msgs = append(j.msgs, err.Error())
	j.Status = statusError
}
