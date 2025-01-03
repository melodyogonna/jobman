package jobman

import (
	"context"
	"log"
	"time"
)

type defaultPooler struct {
	s Backend
}

func (pooler *defaultPooler) Poll(p JobPool) {
	for {
		time.Sleep(time.Minute)
		log.Print("Looking for jobs")
		jobs, err := pooler.s.FindDue(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Found %d jobs", len(jobs))
		for _, job := range jobs {
			p <- job
		}
	}
}
