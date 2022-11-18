package awsinstances

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
)

type promise struct {
	waiter  *ec2.InstanceRunningWaiter
	ref     providers.InstanceRef
	manager *instances.Manager
}

func newPromise(
	client *ec2.Client, ref providers.InstanceRef, manager *instances.Manager,
) *promise {
	return &promise{
		waiter:  ec2.NewInstanceRunningWaiter(client),
		ref:     ref,
		manager: manager,
	}
}

func (p *promise) Await(ctx context.Context) (providers.Instance, error) {
	// Wait for startup
	output, err := p.waiter.WaitForOutput(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{p.ref.ProviderID},
	}, time.Hour)
	if err != nil {
		return providers.Instance{}, providers.NewAPIError("failed to wait for startup", err)
	}

	// Extract information
	instance := output.Reservations[0].Instances[0]
	result, err := instanceFromAwsInstance(instance, p.manager, p.ref)
	if err != nil {
		return providers.Instance{}, providers.NewFatalError(
			"failed to parse returned instance", err,
		)
	}
	return result, nil
}
