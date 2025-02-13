package jobman

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/melodyogonna/jobman/internal/migration"
)

var schemaname string = "jobman_l1nsxfvzfj"

type postgresBackend struct {
	dbHandler *pgxpool.Pool
}

func (backend postgresBackend) Save(ctx context.Context, job TimedJob) error {
	query := `INSERT INTO %s.jobman (job_type, due_on, data, opts) VALUES ($1, $2, $3, $4)`
	_, err := backend.dbHandler.Exec(context.Background(), fmt.Sprintf(query, schemaname), job.Type(), job.In(), job.Payload(), job.Options())
	return err
}

func (backend postgresBackend) FindDue(ctx context.Context) ([]TimedJob, error) {
	tx, err := backend.dbHandler.Begin(ctx)
	if err != nil {
		return nil, err
	}
	query := `WITH due_jobs AS (SELECT id, job_type, due_on, data FROM %s.jobman WHERE due_on <= now() AND job_status = 'PENDING' FOR UPDATE SKIP LOCKED)
	          UPDATE %s.jobman SET job_status='RUNNING' WHERE id IN (SELECT id FROM due_jobs) RETURNING *
	`
	jobs := make([]TimedJob, 0)
	rows, err := backend.dbHandler.Query(context.Background(), fmt.Sprintf(query, schemaname, schemaname))
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	tx.Commit(ctx)
	defer rows.Close()
	for rows.Next() {
		timedjob := GenericTimedJob{}
		var payload json.RawMessage
		err := rows.Scan(&timedjob.Id, &timedjob.JobType, &timedjob.When, &payload)
		if err != nil {
			return nil, err
		}
		timedjob.Data = payload
		jobs = append(jobs, timedjob)
	}
	return jobs, nil
}

func (backend postgresBackend) MarkComplete(ctx context.Context, job TimedJob) error {
	query := `UPDATE %s.jobman SET completed_on=current_timestamp, job_status='FINISHED' WHERE id = $1`
	_, err := backend.dbHandler.Exec(context.Background(), fmt.Sprintf(query, schemaname), job.ID())
	return err
}

func (backend postgresBackend) setup() {
	err := migration.Init(backend.dbHandler.Config().ConnString())
	if err != nil {
		log.Fatal(err)
	}
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
