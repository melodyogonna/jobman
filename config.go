package jobman

type SetupConfig struct {
	Backend    Backend
	Pooler     *Pooler
	WorkerSize uint
}

type JobOptions struct {
	Retry bool
}

type RetryOptions struct{}
