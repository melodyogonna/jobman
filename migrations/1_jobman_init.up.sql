CREATE SCHEMA IF NOT EXISTS jobman_l1nsxfvzfj;
  CREATE TYPE jobman_l1nsxfvzfj.jobman_status AS ENUM('PENDING', 'RUNNING', 'FINISHED');
  CREATE TABLE IF NOT EXISTS jobman_l1nsxfvzfj.jobman(
    id serial PRIMARY KEY,
    job_type VARCHAR(255) NOT NULL,
    data JSONB,
    opts JSONB,
    due_on TIMESTAMP NOT NULL,
    completed_on TIMESTAMP,
    job_status jobman_l1nsxfvzfj.jobman_status DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  );
