package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// TODO: maybe interface Job and package this

// Queue represents a job queue.
type Queue struct {
	jobs     []*Job
	jobsCh   chan *Job
	nWorkers int
	mux      *sync.Mutex
}

// NewQueue creates a new job queue and launches the specified workers.
func NewQueue(workers int) *Queue {
	queue := Queue{
		jobsCh:   make(chan *Job),
		mux:      &sync.Mutex{},
		nWorkers: workers,
	}

	for i := 0; i < workers; i++ {
		go queue.worker()
	}

	posTicker := time.NewTicker(10 * time.Second)
	go func() {
		for range posTicker.C {
			queue.posMsg()
		}
	}()

	return &queue
}

// posMsg informs all jobs in the queue with a message of their position.
func (q *Queue) posMsg() {
	for _, j := range q.jobs {
		pos := q.GetJobPosition(j) - q.nWorkers + 1
		if pos > 0 {
			t := time.Now().Format("15:04:05-0700")
			m := fmt.Sprintf(t+" "+"Awaiting in queue. Position #%d", pos)
			j.msgs = append(j.msgs, m)
		}
	}
}

func (q *Queue) worker() {
	for j := range q.jobsCh {
		j.Process()
		q.Delete(j)
		q.posMsg()
	}
}

// Length returns the number of pending jobs currently in the queue.
func (q *Queue) Length() int {
	return len(q.jobs)
}

// Add adds a new job at the end of the queue.
func (q *Queue) Add(job *Job) {
	q.mux.Lock()
	q.jobs = append(q.jobs, job)
	go func() { q.jobsCh <- job }()
	q.mux.Unlock()
}

// Delete removes a given job from the queue.
func (q *Queue) Delete(job *Job) {
	q.mux.Lock()
	pos := q.GetJobPosition(job)
	q.jobs = append(q.jobs[:pos], q.jobs[pos+1:]...)
	q.mux.Unlock()
}

// GetJobPosition returns the queue position of a given job.
func (q *Queue) GetJobPosition(job *Job) int {
	for i, j := range q.jobs {
		if j == job {
			return i
		}
	}
	return -1
}

// GetJob returns a job in the queue, given an ID.
func (q *Queue) GetJob(id string) (*Job, error) {
	for _, j := range q.jobs {
		if j.ID == id {
			return j, nil
		}
	}
	return nil, errors.New("not found")
}
