package worker

import (
	"context"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"sync"
	"time"
)

type JobType string

const (
	JobTypeTaskCreated JobType = "task_created"
	JobTypeTaskUpdated JobType = "task_updated"
	JobTypeTaskDeleted JobType = "task_deleted"
)

type Job struct {
	Type       JobType
	TaskID     uuid.UUID
	AssigneeID *uuid.UUID
}

type Worker struct {
	jobs chan Job
	once sync.Once
}

func NewWorker(buffer int) *Worker {
	return &Worker{
		jobs: make(chan Job, buffer),
	}
}

func (w *Worker) Start(ctx context.Context) {
	log.Printf("worker started")

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker stopped")
			return
		case job, ok := <-w.jobs:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("worker panic recovered: %v", r)
					}
				}()
				if !ok {
					log.Println("worker channel closed")
					return
				}
				w.handleJob(job)
			}()
		}
	}
}

func (w *Worker) Enqueue(job Job) {
	select {
	case w.jobs <- job:
	default:
		log.Printf("worker queue is full, job dropped: type=%s task_id=%s", job.Type, job.TaskID)
	}
}

func (w *Worker) handleJob(job Job) {
	log.Printf("job started: type=%s, task_id=%s", job.Type, job.TaskID)
	switch job.Type {
	case JobTypeTaskCreated:
		log.Printf("start processing task_created: task_id=%s", job.TaskID)

		time.Sleep(2 * time.Second)

		if job.AssigneeID != nil {
			log.Printf("notification sent to assignee=%s for task=%s", *job.AssigneeID, job.TaskID)
		} else {
			log.Printf("notification no sent, then assignee no pined for task=%s", job.TaskID)
		}

	case JobTypeTaskUpdated:
		log.Printf("start processing task_updated: task_id=%s", job.TaskID)
		time.Sleep(1 * time.Second)
		log.Printf("post-processing finished for updated task=%s", job.TaskID)

	case JobTypeTaskDeleted:
		log.Printf("start processing task_deleted: task_id=%s", job.TaskID)
		time.Sleep(1 * time.Second)
		log.Printf("cleanup finished for deleted task=%s", job.TaskID)

	default:
		log.Printf("unknown job type: %s", job.Type)
	}
	log.Printf("job finished: type=%s, task_id=%s", job.Type, job.TaskID)
}

//func (w *Worker) Close() {
//	w.once.Do(func() { close(w.jobs) })
//}
