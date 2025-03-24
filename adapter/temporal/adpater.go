package temporal

import (
	"fmt"
	"go.temporal.io/sdk/client"
	"log"
)

type Adapter struct {
	Client client.Client
}

func New() Adapter {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalf(fmt.Sprintf("Unable to connect to temporal server, %v", err))
	}
	defer c.Close()

	return Adapter{
		Client: c,
	}
}
