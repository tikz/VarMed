package main

import (
	"time"
)

/* TODO: finish this and possibly make standalone package to reuse in other projects/microservices */
type Job struct {
	RequestedAt    time.Time
	FinishedAt     time.Time
	RunningSecs    float64
	RequesterIP    string
	RequesterEmail string
	Protein        *Protein
	Variations     map[int]string
}

func NewJob() {

}
