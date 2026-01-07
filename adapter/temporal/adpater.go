package temporal

import (
	"log"

	"go.temporal.io/sdk/client"
)

type Config struct {
	Namespace           string `koanf:"namespace"`
	TemporalHostPort    string `koanf:"temporal_host_port"`
	StartToCloseTimeout uint64 `koanf:"start_to_close_timeout"`
	InitialInterval     uint64 `koanf:"initial_interval"`
	BackoffCoefficient  uint64 `koanf:"backoff_coefficient"`
	MaximumInterval     uint64 `koanf:"maximum_interval"`
	MaximumAttempts     uint64 `koanf:"maximum_attempts"`
}

type Adapter struct {
	client client.Client
}

func New(config Config) Adapter {
	c, err := client.Dial(client.Options{
		Namespace: config.Namespace,
		HostPort:  config.TemporalHostPort,
	})
	if err != nil {
		log.Fatalf("Unable to connect to temporal server, %v", err)
	}

	return Adapter{
		client: c,
	}
}

func (a Adapter) GetClient() client.Client {
	return a.client
}

func (a Adapter) Shutdown() {
	a.client.Close()
}
