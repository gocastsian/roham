package job

import (
	"github.com/gocastsian/roham/jobapp/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.temporal.io/sdk/testsuite"
)

func Test_Greeting(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	activities := repository.New()
	service := Service{
		Repository: activities,
		Config: Config{
			StartToCloseTimeout: 10,
			InitialInterval:     1,
			BackoffCoefficient:  1,
			MaximumInterval:     100,
			MaximumAttempts:     3,
		},
	}

	env.OnActivity(activities.SayHelloInPersian, mock.Anything, "nima").Return("سلام nima", nil)

	env.ExecuteWorkflow(service.Greeting, "nima")

	var res string
	err := env.GetWorkflowResult(&res)
	assert.NoError(t, err)
	assert.Equal(t, "سلام nima", res)
}
