package main

import (
	"errors"
	"fmt"
	"sync"
)

// TODO: maybe interface{} Job

// Queue represents a job queue.
type Queue struct {
	jobs   []*Job
	jobsCh chan *Job
	mux    *sync.Mutex
}

// NewQueue creates a new job queue and launches the specified workers.
func NewQueue(workers int) *Queue {
	queue := Queue{
		jobsCh: make(chan *Job),
		mux:    &sync.Mutex{},
	}

	for i := 0; i < workers; i++ {
		go queue.worker()
	}

	go func() {
		for _, j := range queue.jobs {
			select {
			case j.MsgChan <- "caca":
				fmt.Println("sent")
			default:
				fmt.Println("not")
			}
		}
	}()

	return &queue
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

func (q *Queue) worker() {
	for j := range q.jobsCh {
		fmt.Println("processing", j)
		j.Process()
		q.Delete(j)
	}
}
