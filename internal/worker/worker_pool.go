package worker

import (
	"context"
	"fmt"
	"go_job_queue/internal/job"
	"sync"
)

type Handler func(ctx context.Context, j job.Job) error
type Registry map[string]Handler
type OnCompletion func(j job.Job, err error)

type Pool struct {
	n            int
	jobs         <-chan job.Job
	handlers     Registry
	onCompletion OnCompletion
	wg           sync.WaitGroup
}

func NewPool(n int, jobs <-chan job.Job, handlers Registry, onCompletion OnCompletion) *Pool {
	return &Pool{n: n, jobs: jobs, handlers: handlers, onCompletion: onCompletion}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for j := range p.jobs {
				p.onCompletion(j, p.execute(ctx, j))
			}
		}()
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) execute(parent context.Context, j job.Job) error {
	h, ok := p.handlers[j.Type]
	if !ok {
		return fmt.Errorf("No handler for this job type: %q , id: %q", j.Type, j.ID)
	}

	var ctx context.Context
	var cancel context.CancelFunc

	if j.Timeout > 0 {
		ctx, cancel = context.WithTimeout(parent, j.Timeout)
	} else {
		ctx, cancel = context.WithCancel(parent)
	}
	defer cancel()

	done := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("Handler panic: %v", r)
			}
		}()
		done <- h(ctx, j)
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}

}
