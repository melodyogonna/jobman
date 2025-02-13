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

	Options() JobOptions
}

// TimedJob is a job that shouldn't be worked on immediately. Only Timed jobs will be
// persisted in an external storage.
type TimedJob interface {
	Job

	// When job should be worked on
	In() time.Time

	// Mark job as complete after it has been handled. Persistent jobs needs a way to be marked as complete
	MarkCompleted() error

	// Timed jobs require an identity.
	ID() any
}

// Poller searches some job storage interface internal to it and forwards every job it finds to the job pool
type Poller interface {
	Poll(p JobPool)
}

var storage Backend

// handler is a function responsible for processing due jobs.
// TODO: Require handlers to take a context object where errors can be reported
type handler func(job Job) error

type JobPool chan Job

var jobPool JobPool

var jobHandlers map[string][]handler = make(map[string][]handler)

// RegisterHandlers registers a function that should be called when job of jobType needs to be processed.
func RegisterHandlers(jobType string, handlers ...handler) {
	for _, h := range handlers {

		if handlerExistsForJob(jobType, h) {
			continue
		}

		jobHandlers[jobType] = append(jobHandlers[jobType], h)
	}
}

func handlerExistsForJob(job string, h handler) bool {
	handlers, ok := jobHandlers[job]
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
	log.Printf("Got new job. type: %s", job.Type())
	j := newJob{t: newjobtype, data: job}
	jobpool := getJobPool()
	jobpool <- j
}

// TODO: Determine what to do with errors
func initWorkers(p JobPool, workerSize int) {
	for i := 0; i < workerSize; i++ {
		go worker(p)
	}
}

func initializeNewJobHandlers() {
	RegisterHandlers(newjobtype, newJobHandler)
	RegisterHandlers(newjobtype, newTimedJobHandler)
}

func getJobPool() JobPool {
	if jobPool != nil {
		return jobPool
	}

	jobPool = make(JobPool, 5)
	return jobPool
}

// Init sets up Jobman with default options
// namely - 5 workers, no poller, and no backend.
// TimedJob is disabled when jobman is initialized with default options. Sending timed jobs
// will cause a panic.
func Init() {
	jpool := getJobPool()
	initWorkers(jpool, 5)
	initializeNewJobHandlers()
}

// InitWithOptions sets up Jobman with custom options.
// If you provide a backend but no pooler then jobman will use the default poller, which polls the storage every minute.
func InitWithOptions(options SetupConfig) {
	if options.Backend == nil {
		log.Fatal("Backend must be provided when initializing job with options")
	}
	storage = options.Backend
	jpool := getJobPool()
	if options.WorkerSize > 0 {
		initWorkers(jpool, int(options.WorkerSize))
	} else {
		initWorkers(jpool, 5)
	}
	initializeNewJobHandlers()
	if options.Poller != nil {
		pooler := options.Poller
		go pooler.Poll(jpool)
	} else {
		pooler := defaultPooler{options.Backend}
		go pooler.Poll(jpool)
	}
}
