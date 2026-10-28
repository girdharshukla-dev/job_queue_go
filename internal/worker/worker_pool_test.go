package worker

import (
	"context"
	"fmt"
	"go_job_queue/internal/job"
	"testing"
	"time"
)


func TestPoolExecutesJob(t *testing.T){
	jobs := make(chan job.Job, 1)
	completed := make(chan error, 1)
	handlers := Registry{
		"test": func(ctx context.Context, j job.Job) error {
			fmt.Println("Sending an email, this is a test")
			return nil
		},
	}

	pool := NewPool(1, jobs, handlers, func(j job.Job, err error){
		completed <- err 
	})

	jobs <- job.Job{ID: "1", Type: "test"}
	close(jobs)

	pool.Start(context.Background())
	pool.Wait()

	select {
	case err := <-completed:
		if err != nil {
			t.Fatalf("expected nil error: %v", err)
		}
	case <- time.After(time.Second):
		t.Fatal("completion callback not called")
	}

}

