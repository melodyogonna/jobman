package jobman_test

import (
	"context"
	"testing"
	"time"

	"github.com/melodyogonna/jobman"
)

var dbURL = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

func TestJobIsSaved(t *testing.T) {
	b := jobman.PostgresBackend(dbURL)
	job := jobman.GenericTimedJob{JobType: "Test", When: time.Now().Add(time.Hour * 2), Data: struct{ Num int }{Num: 10}}
	err := b.Save(context.Background(), job)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDueJobsFound(t *testing.T) {
	b := jobman.PostgresBackend(dbURL)
	job := jobman.GenericTimedJob{JobType: "Test", When: time.Now()}
	err := b.Save(context.Background(), job)
	if err != nil {
		t.Fatal(err)
	}

	due, err := b.FindDue(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(due) < 1 {
		t.Error("Due jobs not found")
	}
}
