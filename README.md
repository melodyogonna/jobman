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
  }) // Initialize jobman with default configurations
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
