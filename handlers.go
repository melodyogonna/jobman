package jobman

import (
	"context"
	"errors"
	"log"
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
	if storage == nil {
		log.Fatal("Jobman is not initialized with a backend. Please configure jobman with a backend to enable timed jobs.")
	}

	// Add job to the pool if it is already ready to be worked on
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
	err := storage.Save(context.TODO(), job)
	if err != nil {
		return err
	}
	return nil
}
