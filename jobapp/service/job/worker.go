package job

import (
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

type Worker struct {
	worker worker.Worker
}

func New(client client.Client, queueName string) Worker {
	w := worker.New(client, queueName, worker.Options{})

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalf("worker.Run err: %v", err)
	}

	return Worker{}
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
