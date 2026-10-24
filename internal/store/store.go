package store

import "go_job_queue/internal/job"

type Store interface {
	Get(id string) (job.Job, error)
	Save(j job.Job) error
	ListByStatus(status job.Status) ([]job.Job, error)
	Close() error
}


