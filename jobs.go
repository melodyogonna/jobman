package jobman

import (
	"context"
	"time"
)

var NEWJOBTYPE = "NEWJOB"

// newJob is a private job. It is added to the job pool when a new job is registered, it'll be forwarded to new job handlers.
type newJob struct {
	t    string
	data Job
}

func (job newJob) Type() string {
	return job.t
}

func (job newJob) Payload() any {
	return job.data
}

func (job newJob) Options() JobOptions {
	return JobOptions{}
}

type GenericJob struct {
	JobType string
	Data    any
	Opts    *JobOptions
}

func (job GenericJob) Type() string {
	return job.JobType
}

func (job GenericJob) Payload() any {
	return job.Data
}

func (job GenericJob) Options() *JobOptions {
	return job.Opts
}

// GenericTimedJob is a default job that shouldn't get handled immediately. We need to enforce persistence when dealing with Timed jobs, so
// It requires a way to mark it as complete.
type GenericTimedJob struct {
	JobType string
	Data    any // This should be JSON serializable if we intend to save this in DB
	Opts    *JobOptions
	id      int
	When    time.Time
}

func (job GenericTimedJob) Type() string {
	return job.JobType
}

func (job GenericTimedJob) Payload() any {
	return job.Data
}

func (job GenericTimedJob) Options() JobOptions {
	if job.Opts != nil {
		return *job.Opts
	}
	return JobOptions{}
}

func (job GenericTimedJob) MarkCompleted() error {
	err := storage.MarkComplete(context.TODO(), &job)
	return err
}

func (job GenericTimedJob) In() time.Time {
	return job.When
}

func (job GenericTimedJob) ID() any {
	return job.id
}
