package jobman

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Backend defines an interface for persisting jobs to an external persistent storage.
type Backend interface {
	// Save job, return an error job could not be saved
	Save(job *TimedJob) error
	// Find due jobs.
	FindDue() ([]TimedJob, error)
	// Mark a job as complete
	MarkComplete(job TimedJob) error
}

type postgresBackend struct {
	dbHandler *pgxpool.Pool
}

func (backend postgresBackend) Save(job *TimedJob) error {
	d := *job
	query := `INSERT INTO jobman (job_type, due_on, data) VALUES ($1, $2, $3)`
	_, err := backend.dbHandler.Exec(context.Background(), query, d.Type(), d.In(), d.Payload())
	return err
}

func (backend postgresBackend) FindDue() ([]TimedJob, error) {
	ctx := context.Background()
	tx, err := backend.dbHandler.Begin(ctx)
	if err != nil {
		return nil, err
	}
	query := `WITH due_jobs AS (SELECT id, job_type, due_on, data FROM jobman WHERE due_on <= now() AND status = 'PENDING' FOR UPDATE SKIP LOCKED)
	  UPDATE due_jobs SET status='RUNNING' RETURNING *
	`
	jobs := make([]TimedJob, 0)
	rows, err := backend.dbHandler.Query(context.Background(), query)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	tx.Commit(ctx)
	defer rows.Close()
	for rows.Next() {
		timedjob := GenericTimedJob{}
		var payload *any
		err := rows.Scan(&timedjob.id, &timedjob.JobType, &timedjob.When, &payload)
		if err != nil {
			return nil, err
		}
		timedjob.Data = *payload
		jobs = append(jobs, timedjob)
	}
	return jobs, nil
}

func (backend postgresBackend) MarkComplete(job TimedJob) error {
	query := `UPDATE jobman SET completed_on=current_timestamp WHERE id = $1`
	_, err := backend.dbHandler.Exec(context.Background(), query, job.ID())
	return err
}

func (backend postgresBackend) setup() {
}

func PostgresBackend(url string) *postgresBackend {
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal(err)
	}
	backend := &postgresBackend{db}
	backend.setup()
	return backend
}
