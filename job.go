package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"respdb/pdb"
	"respdb/uniprot"
	"sort"
	"strconv"
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
	Name      string    `json:"name"`
	UniProtID string    `json:"uniprotId"`
	PDBIDs    []string  `json:"pdbIds"`
	ClinVar   bool      `json:"clinvar"`
	SAS       []string  `json:"sas"`
	IP        string    `json:"ip"`
	Email     string    `json:"email"`
	Time      time.Time `json:"time"`
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

// generateID returns a SHA256 hash of UniProtID+sorted PDBIDs+sorted SASs.
func generateID(r *JobRequest) string {
	unpID := []byte(r.UniProtID)

	pdbIDs := r.PDBIDs
	sort.Strings(pdbIDs)
	pdbBytes := []byte(strings.Join(pdbIDs, ""))

	sas := r.SAS
	sort.Strings(sas)
	sasBytes := []byte(strings.Join(sas, ""))

	b := bytes.Join([][]byte{unpID, pdbBytes, sasBytes}, []byte(""))
	hash := sha256.Sum256(b)

	return hex.EncodeToString(hash[:])
}

// NewJob returns a new job instance.
func NewJob(request *JobRequest) *Job {
	j := &Job{Request: request}
	j.ID = generateID(request)

	return j
}

// Process runs the pipeline for the job.
func (j *Job) Process(cli bool) {
	j.Status = statusProcess
	j.Started = time.Now()

	unp, err := loadUniProt(j.Request.UniProtID)
	if err != nil {
		j.fail(err)
		return
	}

	substs, err := loadSAS(unp.Sequence, j.Request.SAS)
	if err != nil {
		j.fail(fmt.Errorf("parse SAS list: %v", err))
		return
	}

	msgChan := make(chan string, 100)
	j.Pipeline, _ = NewPipeline(unp, j.Request.PDBIDs, substs, msgChan)

	go func() {
		for m := range msgChan {
			if cli {
				fmt.Println(m)
			} else {
				j.msgs = append(j.msgs, m)
			}
		}
	}()

	err = j.Pipeline.Run()
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
	log.Printf("error %s %s: %v", j.Request.UniProtID, j.Request.PDBIDs, err)
	j.msgs = append(j.msgs, err.Error())
	j.Error = err
	j.Status = statusError
}

// loadSAS parses and validates a slice of formatted SAS strings.
func loadSAS(seq string, sas []string) ([]*uniprot.SAS, error) {
	var parsedSAS []*uniprot.SAS

	for _, s := range sas {
		r, _ := regexp.Compile("(.)([0-9]*)(.)")
		m := r.FindStringSubmatch(s)
		if m == nil {
			return nil, errors.New("bad SAS format:" + s)
		}

		from := m[1]
		pos, _ := strconv.ParseInt(m[2], 10, 64)
		to := m[3]

		if pos <= 0 {
			return nil, errors.New(s + " position must be 1 or greater")
		}

		if !pdb.IsAminoacid(from) {
			return nil, errors.New(s + " not an aminoacid: " + from)
		}
		if !pdb.IsAminoacid(to) {
			return nil, errors.New(s + " not an aminoacid: " + to)
		}

		if from == to {
			return nil, errors.New(s + " has same aminoacids, not a SAS")
		}

		unpAa := string(seq[pos-1])
		if from != unpAa {
			errStr := fmt.Sprintf("SAS %s: position %d in UniProt seq has Aa %s, not %s",
				s, pos, unpAa, from)
			return nil, errors.New(errStr)
		}

		parsedSAS = append(parsedSAS, &uniprot.SAS{
			Position: pos,
			FromAa:   from,
			ToAa:     to,
		})
	}

	return parsedSAS, nil
}
