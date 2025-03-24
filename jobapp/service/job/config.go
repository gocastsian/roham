package job

type Config struct {
	GreetingQueueName   string `koanf:"greeting_queue_name"`
	StartToCloseTimeout uint64 `koanf:"start_to_close_timeout"`
	InitialInterval     uint64 `koanf:"initial_interval"`
	BackoffCoefficient  uint64 `koanf:"backoff_coefficient"`
	MaximumInterval     uint64 `koanf:"maximum_interval"`
	MaximumAttempts     uint64 `koanf:"maximum_attempts"`
}
