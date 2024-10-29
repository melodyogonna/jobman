package jobman

import (
	"log"
	"time"
)

type defaultPooler struct {
	s Backend
}

func (pooler *defaultPooler) Pool(p JobPool) {
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
