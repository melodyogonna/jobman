package jobman

import "time"

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

// GenericJob is a default without any strong typing. It accepts all values and is scheduled to get handled immediately.
type GenericJob struct {
	JobType string
	Data    any // This should be JSON serializable if we intend to save this in DB
}

func (job GenericJob) Type() string {
	return job.JobType
}

func (job GenericJob) Payload() any {
	return job.Data
}

// GenericTimedJob is a default job that shouldn't get handled immediately. We need to enforce persistence when dealing with Timed jobs, so
// It requires a way to mark it as complete.
type GenericTimedJob struct {
	GenericJob
	Id   int
	When time.Time
}

func (job GenericTimedJob) MarkCompleted() error {
	err := dbHandler.MarkComplete(&job)
	return err
}

func (job GenericTimedJob) In() time.Time {
	return job.When
}

func (job GenericTimedJob) ID() int {
	return job.Id
}
