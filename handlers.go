package jobman

import (
	"errors"
	"time"
)

func newJobHandler(job Job) error {
	newJob, ok := job.(newJob)
	if !ok {
		return errors.New("Not a new job")
	}

	// Determine whether this job is a timed job. Don't proceed if it is, we have a timedJob handler
	if _, ok := newJob.Payload().(TimedJob); ok {
		return nil
	}

	jobPool := getJobPool()
	jobPool <- newJob.data

	return nil
}

func newTimedJobHandler(job Job) error {
	newJob, ok := job.(newJob)
	if !ok {
		return errors.New("Not a new job")
	}

	// Determine whether this job is a timed job. Don't proceed if it is not
	timedJob, ok := newJob.Payload().(TimedJob)
	if !ok {
		return nil
	}

	jobDue := timedJob.In()
	now := time.Now()
	if !jobDue.After(now) {
		jobPool := getJobPool()
		jobPool <- timedJob
		return nil
	}

	err := saveJob(timedJob)
	if err != nil {
		return err
	}
	return nil
}

func saveJob(job TimedJob) error {
	err := dbHandler.Save(&job)
	if err != nil {
		return err
	}
	return nil
}
