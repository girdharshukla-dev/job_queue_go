package store

import (
	"go_job_queue/internal/job"
	"testing"
	"time"
)

func TestSaveGetListReopen(t *testing.T) {
	dir := t.TempDir()

	s, err := NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	now := time.Now()
	a := job.Job{
		ID:          "a",
		Type:        "email",
		Payload:     []byte("hello"),
		Status:      job.StatusPending,
		MaxAttempts: 3,
		RunAt:       now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	b := a
	b.ID = "b"
	b.Status = job.StatusCompleted

	err = s.Save(a)
	if err != nil {
		t.Fatalf("save a: %v", err)
	}

	err = s.Save(b)
	if err != nil {
		t.Fatalf("save b: %v", err)
	}

	got, err := s.Get("a")
	if err != nil {
		t.Fatalf("get a: %v", err)
	}
	checkJob(t, got, a)

	pending, err := s.ListByStatus(job.StatusPending)
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(pending) != 1 || pending[0].ID != "a" {
		t.Fatalf("list pending: %+v", err)
	}

	err = s.Close()
	if err != nil {
		t.Fatalf("Error in closing: %v", err)
	}

	s2, err := NewBadgerStore(dir)
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	defer s2.Close()
	got, err = s2.Get("b")
	if err != nil {
		t.Fatalf("Error in get b: %v", err)
	}
	checkJob(t, got, b)

}

func checkJob(t *testing.T, got, want job.Job) {
	t.Helper()
	if got.ID != want.ID || got.Type != want.Type || got.Status != want.Status || got.MaxAttempts != want.MaxAttempts ||
		string(got.Payload) != string(want.Payload) || !got.RunAt.Equal(want.RunAt) || !got.CreatedAt.Equal(want.CreatedAt) {
		t.Fatalf("job mismatch:\n got: %+v\n want: %+v", got, want)
	}
}
