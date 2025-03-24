package job

import (
	"fmt"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

type Workflow struct {
}

func Greeting(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumInterval:    time.Second * 5,
			MaximumAttempts:    10,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var persianRes string
	err := workflow.ExecuteActivity(ctx, SayHelloInPersian, name).Get(ctx, &persianRes)
	if err != nil {
		return "", fmt.Errorf("greeting in persain failed")
	}

	return persianRes, err
}
