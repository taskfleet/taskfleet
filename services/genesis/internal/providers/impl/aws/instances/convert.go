package awsinstances

import (
	"fmt"
	"net"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
)

func instanceFromAwsInstance(
	instance types.Instance, manager *instances.Manager, ref providers.InstanceRef,
) (providers.Instance, error) {
	// Assemble spec
	instanceType, err := manager.Get(string(instance.InstanceType))
	if err != nil {
		return providers.Instance{}, fmt.Errorf("failed to get instance type: %s", err)
	}
	spec := providers.InstanceSpec{
		InstanceType: instanceType,
		IsSpot:       instance.InstanceLifecycle == types.InstanceLifecycleTypeSpot,
	}

	// Assemble status
	internalIP := net.ParseIP(*instance.PrivateIpAddress)
	if internalIP == nil {
		return providers.Instance{}, fmt.Errorf("encountered invalid internal IP")
	}
	externalIP := net.ParseIP(*instance.PublicIpAddress)
	if externalIP == nil {
		return providers.Instance{}, fmt.Errorf("encountered invalid external IP")
	}

	status := providers.InstanceStatus{
		CreationTimestamp: *instance.LaunchTime,
		Network: providers.InstanceNetworkStatus{
			InternalIP:       internalIP,
			ExternalIP:       externalIP,
			InternalHostname: *instance.PrivateDnsName,
		},
	}

	// Return instance
	return providers.Instance{Ref: ref, Spec: spec, Status: status}, nil
}
