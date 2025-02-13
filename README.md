# jobman

Easy and powerful background jobs for Go applications.

# install

```sh
go get -u github.com/melodyogonna/jobman
```

# Usage

You can start using Jobman in two ways.

## Using jobman with default options

The easiest way to start using jobman is by Initializing it without any options.
This initializes jobman without any storage support, and with 5 workers.

```go
package main

import (
  "github.com/melodyogonna/jobman"
  "log"
)

type EmailData struct {
      email string
      templateId string
}

func HandleEmailJob(job jobman.Job) error {
  log.Print("Email job handler called")
  data = job.Payload().(EmailData)
  // ...handle sending email
  return nil
}

func main(){
  jobman.Init() // Initialize jobman with default configurations
  jobman.RegisterHandlers("sendEmail", HandleEmailJob)

  job := jobman.GenericJob{
    JobType: "sendEmail",
    Data: EmailData{email:"johndoe@email.com", templateId: "testEmailTemplateId"}
  }
  jobman.WorkOn(job)
}
```

When a job is added to Jobman's job pool it'll be logged. For our example above, something like this is expected:

```sh
2025/02/13 15:19:15 worker: f222b4ab-1676-4245-bddd-30f06b902234 - handling job with type: NEWJOB                                             [0/8052]
2025/02/13 15:19:15 2 handlers registered for job type: NEWJOB. Forwarding ...
2025/02/13 15:19:15 worker: bddfa149-ab61-4dcb-a100-3b023be30996 - handling job with type: sendEmail
2025/02/13 15:19:15 1 handlers registered for job type: sendEmail. Forwarding ...
```

## Using jobman with custom options

You can configure Jobman to use custom storage backend - this allows you to create timed jobs. Timed jobs are associated with a future time
when the job should be handled.

```go

package main

import (
  "github.com/melodyogonna/jobman"
  "log"
  "os"
  "time"
)

func HandleEmailJob(job jobman.Job) error {
  log.Print("Email job handler called")
  // ...handle sending email
  return nil
}

func main(){
  jobman.InitWithOptions(jobman.SetupConfig{
    Backend: jobman.PostgresBackend(os.Getenv("DATABASE_URL")),
    WorkerSize: 10
  }) // Initialize jobman with custom configurations
  jobman.RegisterHandlers("sendEmail", HandleEmailJob)

  job := jobman.GenericTimedJob{
    JobType: "sendEmail",
    When: time.Now().Add(time.HOUR * 24)
    Data: EmailData{email:"johndoe@email.com", templateId: "testEmailTemplateId"}
  }
  jobman.WorkOn(job)
}
```

Jobman will panic if you tried to make it work on a timed job without setting up a backend.

### Timed Jobs

Timed jobs have timing attached, jobman will save these jobs using the specified backend.

# Components

My goal with this library is to make something where simple parts compose together intuitively. To that end, Jobman's core has 3 composing parts:

1. Poller
2. Backend
3. Job

## Poller

The purpose of a poller is to check some external source for due jobs, then add any jobs found to a job pool. A poller has the interface:

```go

type JobPool chan Job

type Poller interface {
	Poll(p JobPool)
}
```

The default poller checks whatever backend passed during initialization every minute. You can create a custom poller and pass it to Jobman during initialization:

```go
package main

import (
  "github.com/melodyogonna/jobman"
  "time"
)

// customPoller checks for jobs every hour
type customPoller struct {
  storage jobman.Backend
}
func (p CustomPoller) Poll(pool jobman.JobPool) {
  for {
    time.Sleep(time.Hour)
    due, err := p.storage.FindDue()
    if err != nil {
      // do error handling
      return
    }
    for _, job := range due {
      pool <- job
    }
  }
}

func GetCustomPoller(backend jobman.Backend) jobman.Pooler {
  return &customPoller{storage: backend}
}

func main(){
  jobman.InitWithOptions(jobman.SetupConfig{
    Backend: jobman.PostgresBackend(os.Getenv("DATABASE_URL")),
    WorkerSize: 10,
    Poller: GetCustomPoller(jobman.PostgresBackend(os.Getenv("DATABASE_URL")))
  }) // Initialize jobman with custom configurations
}
```

Jobman will use the default 1-minute poller if you don't pass any during setup.

# Roadmap

Some ideas about the things I intend to implement in the future

- [ ] Job retries
- [ ] Redis backend
- [ ] Nats backend
- [ ] Action hooks - To monitor different states of a job
- [ ] More tests
