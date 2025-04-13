package temporal

import (
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	worker worker.Worker
}

func NewWorker(client client.Client, queueName string, options worker.Options) Worker {
	w := worker.New(client, queueName, options)

	return Worker{worker: w}
}

func (w Worker) RegisterWorkflow(wf interface{}) {
	w.worker.RegisterWorkflow(wf)
}

func (w Worker) RegisterActivity(act interface{}) {
	w.worker.RegisterActivity(act)
}

func (w Worker) Start() error {
	err := w.worker.Start()
	if err != nil {
		return fmt.Errorf("can not start worker err: %v", err)
	}

	return nil
}

func (w Worker) Stop() {
	w.worker.Stop()
}
