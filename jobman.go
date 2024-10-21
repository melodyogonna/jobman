package jobman

import (
	"fmt"
	"log"
	"time"
)

// Job defines a standard job interface with properties a job should have.
type Job interface {
	// Retrieve job type
	Type() string
	// Retrieve Job data
	Payload() any
}

// TimedJob is a job that shouldn't be worked on immediately
type TimedJob interface {
	Job

	// When job should be worked on
	In() time.Time

	// Mark job as complete after it has been handled. Persistent jobs needs a way to be marked as complete
	MarkCompleted() error

	ID() int
}

type JobResult interface {
	Result() (any, error)
}

// Pooler searches some job storage interface internal to it and forwards every job it finds to the job pool
type Pooler interface {
	Pool(p pool)
}

type storage interface {
	Save(job *TimedJob) error
	FindDue() ([]TimedJob, error)
	MarkComplete(job TimedJob) error
}

var dbHandler storage

// handler is a function responsible for processing due jobs.
// TODO: Require handlers to take a context object where errors can be reported
type handler func(job Job) error

type pool chan Job

var jobPool pool

var handlers map[string][]handler = make(map[string][]handler)

// RegisterHandler registers a function that should be called when job of jobType needs to be processed.
func RegisterHandler(jobType string, handler handler) {
	if handlerExistsForJob(jobType, handler) {
		return
	}

	handlers[jobType] = append(handlers[jobType], handler)
}

func handlerExistsForJob(job string, h handler) bool {
	handlers, ok := handlers[job]
	if !ok {
		return false
	}

	for _, handler := range handlers {
		addr1 := fmt.Sprintf("%v", handler)
		addr2 := fmt.Sprintf("%v", h)
		if addr1 == addr2 {
			return true
		}
	}

	return false
}

// WorkOn gives Jobman a job to work on. If the Job is a TimedJob it'll be saved to the database, otherwise it'll be handled immediately
func WorkOn(job Job) {
	j := newJob{t: NEWJOBTYPE, data: job}
	jobpool := getJobPool()
	jobpool <- j
}

// TODO: Determine what to do with errors
func initWorkers(p pool) {
	workerSize := 5
	for i := 0; i < workerSize; i++ {
		go worker(p)
	}
}

func initializeNewJobHandlers() {
	RegisterHandler(NEWJOBTYPE, newJobHandler)
	RegisterHandler(NEWJOBTYPE, newTimedJobHandler)
}

func getJobPool() pool {
	if jobPool != nil {
		return jobPool
	}

	jobPool = make(pool, 5)
	return jobPool
}

func StartWork(s storage) {
	dbHandler = s
	log.Print("Initializing job pool")
	jpool := getJobPool()
	log.Print("Starting workers")
	initWorkers(jpool)
	log.Print("Starting default job handlers")
	initializeNewJobHandlers()
	log.Print("Starting default job pooler")
	pooler := defaultPooler{s}
	go pooler.Pool(jpool)
	log.Print("Job man has been initialized and is now working.")
}
