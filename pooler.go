package jobman

import (
	"log"
	"time"
)

// DefaultPooler pools the storage for jobs every six hours
type defaultPooler struct {
	s storage
}

func (pooler *defaultPooler) Pool(p pool) {
	for {
		time.Sleep(time.Minute)
		log.Print("Looking for jobs")
		jobs, error := pooler.s.FindDue()
		if error != nil {
			return
		}
		log.Printf("Found %d jobs", len(jobs))
		for _, job := range jobs {
			p <- job
		}
	}
}
