package job

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusCompleted Status = "completed"
	StatusDead Status = "dead"
)

const DefaultMaxAttempts int64 = 1

type Job struct {
	ID string
	Type string
	Payload []byte
	Status Status
	CurrentAttempts int
	MaxAttempts int
	RunAt time.Time
	Timeout time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
	Error string
}

func New(j Job) (Job, error) {
	
	if j.Type == "" || j.MaxAttempts < 0 || time.Now().After(j.RunAt) || j.Timeout <= 0{
		return Job{}, errors.New("Invalid fields")
	}

	if j.MaxAttempts == 0 {
		j.MaxAttempts = int(DefaultMaxAttempts)
	}
	
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	j.Status = StatusPending
	j.CurrentAttempts = 0

	if j.MaxAttempts == 0 {
		j.MaxAttempts = 1
	}
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()
	return j, nil
}

func UpdateStatus(j *Job, s Status) error {
	if j.Status == StatusCompleted || j.Status == StatusDead {
		return errors.New("Job already dead or completed")
	}

	j.Status = s
	j.UpdatedAt = time.Now()
	return nil
}


func Encode(j Job) ([]byte, error) {
	encoded, err := json.Marshal(j)

	if err != nil {
		return []byte{}, errors.New("Error in json.Marshal in Decode")
	}

	return encoded, nil
}

func Decode(bytes []byte) (Job, error) {
	var j Job

	err := json.Unmarshal(bytes, &j)
	if err != nil {
		return Job{}, errors.New("Error in json.Unmarshal in Decode")
	}
	
	return j, nil
}

