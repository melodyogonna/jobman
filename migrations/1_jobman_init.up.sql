CREATE TYPE jobman_L1nsxfVZfj_status AS ENUM('PENDING', 'RUNNING', 'FINISHED');
CREATE TABLE IF NOT EXISTS jobman_L1nsxfVZfj(
  id serial PRIMARY KEY,
  job_type VARCHAR(255) NOT NULL,
  data JSONB,
  due_on TIMESTAMP NOT NULL,
  completed_on TIMESTAMP,
  job_status jobman_L1nsxfVZfj_status DEFAULT 'PENDING',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
