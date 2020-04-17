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
)

// JobRequest represents a single job request from an user,
// contains the inputs and extra client data.
type JobRequest struct {
	UniProtID     string   `json:"uniprot_id"`
	PDBIDs        []string `json:"pdbs"`
	ClinVar       bool     `json:"clinvar"`
	VariationsPos []int    `json:"variations_pos"`
	VariationsAA  []string `json:"variations_aa"`

	IP    string
	Email string
	Time  time.Time
}

type Job struct {
	ID       string
	Request  *JobRequest
	Pipeline *Pipeline `json:"-"`

	Status  int
	Started time.Time
	Ended   time.Time

	msgChan chan string
	Error   error `json:"-"`
}

// generateID returns a SHA256 hash of UniProtID+joined PDBIDs
func (j *Job) generateID() string { // TODO: include variations in hash after implementing that

	unpID := []byte(j.Request.UniProtID)
	pdbIDs := []byte(strings.Join(j.Request.PDBIDs, ""))
	b := bytes.Join([][]byte{unpID, pdbIDs}, []byte(""))
	hash := sha256.Sum256(b)

	return hex.EncodeToString(hash[:])
}

func NewJob(request *JobRequest) Job {
	j := Job{Request: request}
	j.msgChan = make(chan string, 100)

	j.ID = j.generateID()

	return j
}

func (j *Job) Process() {
	j.Status = statusProcess
	j.Started = time.Now()

	j.Pipeline, _ = NewPipeline(j.Request.UniProtID, j.Request.PDBIDs, j.msgChan)
	err := j.Pipeline.Run()
	if err != nil {
		j.fail(err)
		return
	}
	j.Ended = time.Now()

	j.Status = statusDone

	err = WriteJob(j)
	if err != nil {
		j.fail(err)
		return
	}

	// Either "SUCCESS" or "FAILED" is the specific message that
	// the frontend expects from the WebSocket to proceed.
	j.msgChan <- "SUCCESS"
	close(j.msgChan)
}

func (j *Job) fail(err error) {
	j.msgChan <- err.Error()
	j.msgChan <- "FAILED"
	close(j.msgChan)
}
