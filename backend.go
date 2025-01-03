package jobman

import (
	"context"
)

// Backend defines an interface for persisting jobs to an external persistent storage.
type Backend interface {
	// Save job, return an error job could not be saved
	Save(ctx context.Context, job TimedJob) error
	// Find due jobs.
	FindDue(ctx context.Context) ([]TimedJob, error)
	// Mark a job as complete
	MarkComplete(ctx context.Context, job TimedJob) error
}
