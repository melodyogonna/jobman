package jobman

import (
	"database/sql"
	"log"
)

// Backend defines an interface for persisting jobs to an external persistent storage.
type Backend interface {
	// Save job, return an error job could not be saved
	Save(job *TimedJob) error
	// Find due jobs.
	FindDue() ([]TimedJob, error)
	// Mark a job as complete
	MarkComplete(job *TimedJob) error
}

type postgresBackend struct {
	dbHandler *sql.DB
}

func (backend postgresBackend) Save(job *TimedJob) error {
	d := *job
	query := `INSERT INTO jobman (job_type, due_on, data) VALUES ($1, $2, $3)`
	_, err := backend.dbHandler.Exec(query, d.Type(), d.In(), d.Payload())
	return err
}

func (backend postgresBackend) FindDue() ([]TimedJob, error) {
	query := `SELECT id, job_type, due_on, data FROM jobman WHERE due_on <= now() AND completed_on IS NULL`
	jobs := make([]TimedJob, 0)
	rows, err := backend.dbHandler.Query(query)
	if err != nil {
		return nil, err
	}
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
	_, err := backend.dbHandler.Exec(query, job.ID())
	return err
}

func PostgresBackend(url string) *postgresBackend {
	db, err := sql.Open("psql", url)
	if err != nil {
		log.Fatal(err)
	}
	return &postgresBackend{db}
}
