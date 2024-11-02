# jobman

Easy and powerful background jobs for Go applications.

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

func HandleEmailJob(job jobman.Job) error {
  log.Print("Email job handler called")
  // ...handle sending email
  return nil
}

func main(){
  jobman.Init() // Initialize jobman with default configurations
  jobman.RegisterHandlers("sendEmail", HandleEmailJob)

  job := jobman.GenericJob{
    JobType: "sendEmail",
    Data: struct{
      email string
      templateId string
    }{email:"johndoe@email.com", templateId: "testEmailTemplateId"}
  }
  jobman.WorkOn(job)
}
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
    Data: struct{
      email string
      templateId string
    }{email:"johndoe@email.com", templateId: "testEmailTemplateId"}
  }
  jobman.WorkOn(job)
}
```

Jobman will panic if you tried to make it work on a timed job without setting up a backend.

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
