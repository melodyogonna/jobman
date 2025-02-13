package jobman

type SetupConfig struct {
	Backend    Backend
	Poller     Poller
	WorkerSize uint
}

type JobOptions struct {
	Retry bool `json:"retry"`
}

type RetryOptions struct{}
