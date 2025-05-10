package job

type Event struct {
	WorkflowId   string
	WorkflowName string
	QueueName    string
	Args         map[string]any
}
