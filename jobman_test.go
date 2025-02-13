package jobman

import (
	"testing"
	"time"
)

func simpleTestSetup(t *testing.T) {
	t.Helper()
	Init()
}

func TestHandlerRegisteration(t *testing.T) {
	simpleTestSetup(t)
	f := func(j Job) error { return nil }
	jobType := "TestJob"
	RegisterHandlers(jobType, f)

	if !handlerExistsForJob(jobType, f) {
		t.Fatal("Handler is not registered")
	}
}

func TestJobForwardedToHandler(t *testing.T) {
	simpleTestSetup(t)
	wasCalled := make(chan bool)
	f := func(j Job) error { wasCalled <- true; return nil }
	jobType := "TestJob"
	Init()
	RegisterHandlers(jobType, f)
	job := GenericJob{JobType: jobType}
	WorkOn(job)
	select {
	case <-wasCalled:
		break
	case <-time.After(time.Second * 2):
		t.Fatal("Time out. Handler not called")
	}
}

func BenchmarkGenericJob(b *testing.B) {
	Init()
	f := func(j Job) error { return nil }
	jobType := "TestJob"
	RegisterHandlers(jobType, f)
	job := GenericJob{JobType: jobType, Data: struct{ Num int }{Num: 10}}
	for range b.N {
		WorkOn(job)
	}
}
