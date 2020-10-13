package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// QueueStatus represents a single message about the current queue status.
type QueueStatus struct {
	TotalJobs int              `json:"totalJobs"`
	Jobs      []QueueStatusJob `json:"jobs"`
	MyJobs    []QueueStatusJob `json:"myJobs"`
}

// QueueStatusJob holds simplified and non identifying data about a single job being processed.
type QueueStatusJob struct {
	Position    int     `json:"position"`
	ID          string  `json:"id"`
	ShortID     string  `json:"shortId"`
	Progress    float64 `json:"progress"`
	ProgressPDB float64 `json:"progressPdb"`
	Elapsed     string  `json:"elapsed"`
	PDBs        int     `json:"pdbs"`
	Variants    int     `json:"variants"`
}

// WSJobEndpoint handles WebSocket /ws/job/:jobID
func WSJobEndpoint(c *gin.Context) {
	id := c.Param("jobID")
	queue := c.MustGet("queue").(*Queue)
	job, err := queue.GetJob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	wsJobHandler(c.Writer, c.Request, job)
}

func wsJobHandler(w http.ResponseWriter, r *http.Request, j *Job) {
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

	// Show last 10 messages only when reconnecting
	if i > 10 {
		i = i - 10
	}

	for {
		select {
		case <-msgTicker.C:
			if i < len(j.msgs) {
				msg := j.msgs[i]
				ws.WriteMessage(websocket.TextMessage, []byte(msg))
				i++

				if j.Status == statusDone ||
					j.Status == statusSaved {
					return
				}
			}
		}

	}
}

// WSQueueEndpoint handles WebSocket /ws/queue
func WSQueueEndpoint(c *gin.Context) {
	queue := c.MustGet("queue").(*Queue)
	wsQueueHandler(c.Writer, c.Request, queue, c.ClientIP())
}

func wsQueueHandler(w http.ResponseWriter, r *http.Request, q *Queue, clientIP string) {
	upg := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	upg.CheckOrigin = func(r *http.Request) bool { return true } // TODO: remove in production, unsafe

	ws, err := upg.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade fail: %+v", err)
		return
	}

	msgTicker := time.NewTicker(1000 * time.Millisecond)
	defer func() {
		msgTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case <-msgTicker.C:
			msg, err := json.Marshal(queueStatus(q, clientIP))
			if err == nil {
				ws.WriteMessage(websocket.TextMessage, msg)
			}
		}

	}
}

func queueStatus(q *Queue, clientIP string) (qs QueueStatus) {
	qs.TotalJobs = len(q.jobs)

	for i, job := range q.jobs {
		if i < 2 && job.Pipeline != nil {
			qsJob := QueueStatusJob{
				Position:    i + 1,
				ShortID:     job.ID[:5],
				Elapsed:     time.Now().Sub(job.Started).Truncate(time.Second).String(),
				Progress:    job.Pipeline.Progress,
				ProgressPDB: job.Pipeline.ProgressPDB,
				PDBs:        len(job.Request.PDBIDs),
				Variants:    len(job.Request.Variants),
			}
			qs.Jobs = append(qs.Jobs, qsJob)
		}

		if job.Request.IP == clientIP {
			qs.MyJobs = append(qs.MyJobs, QueueStatusJob{
				Position: i + 1,
				ID:       job.ID,
				ShortID:  job.ID[:5],
				PDBs:     len(job.Request.PDBIDs),
				Variants: len(job.Request.Variants),
			})
		}

	}

	return
}
