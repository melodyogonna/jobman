package jobman

import (
	"testing"
	"time"
)

func TestHandlerRegisteration(t *testing.T) {
	f := func(j Job) error { return nil }
	jobType := "TestJob"
	RegisterHandler(jobType, f)

	if !handlerExistsForJob(jobType, f) {
		t.Fatal("Handler is not registered")
	}
}

func TestJobForwardedToHandler(t *testing.T) {
	wasCalled := make(chan bool)
	f := func(j Job) error { wasCalled <- true; return nil }
	jobType := "TestJob"
	RegisterHandler(jobType, f)
	job := GenericJob{JobType: jobType, Data: struct{ Num int }{Num: 10}}
	WorkOn(job)
	select {
	case <-wasCalled:
		break
	case <-time.After(time.Second * 1):
		t.Fatal("Time out. Handler not called")
	}
}

func BenchmarkGenericJob(b *testing.B) {
	f := func(j Job) error { return nil }
	jobType := "TestJob"
	RegisterHandler(jobType, f)
	job := GenericJob{JobType: jobType, Data: struct{ Num int }{Num: 10}}
	for range b.N {
		WorkOn(job)
	}
}
