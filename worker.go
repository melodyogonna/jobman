package jobman

import (
	"log"

	"github.com/google/uuid"
)

// Worker receives jobs from job pool channel and passes it over to handlers subscribed to such pools.
func worker(jobPool JobPool) {
	workerID := uuid.New().String()
	for {
		job := <-jobPool
		jobType := job.Type()
		log.Printf("worker: %s - handling job with type: %s", workerID, jobType)
		h, ok := jobHandlers[jobType]
		if !ok {
			continue
		}
		for _, handler := range h {
			err := handler(job)
			if err != nil {
				// TODO: Support error handlers that'll take this error
				log.Print(err)
			}
		}
		// mark complete if timedJob job
		if timedJob, ok := job.(TimedJob); ok {
			// TODO: Determine implementing support for retries
			timedJob.MarkCompleted()
		}
	}
}
