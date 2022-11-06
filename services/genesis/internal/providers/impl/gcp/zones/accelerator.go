package gcpzones

import (
	"context"
	"fmt"
	"math/bits"
	"path"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/borchero/zeus/pkg/zeus"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"go.uber.org/zap"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// Accelerator describes a single type of accelerator in a single zone.
type Accelerator struct {
	// The fully qualified name of the accelerator.
	uri string
	// The kind of accelerator.
	kind typedefs.GPUKind
	// The maximum number of accelerators of this type for a particular instance.
	maxCountPerInstance uint16
}

// MaxCount returns the maximum number of accelerators to be attached to an instance.
func (a Accelerator) MaxCount() uint16 {
	return a.maxCountPerInstance
}

// Config returns the accelerator configuration for the provided accelerator count.
func (a Accelerator) Config(count uint16) (*computepb.AcceleratorConfig, error) {
	if count > a.maxCountPerInstance {
		return nil, fmt.Errorf(
			"too many accelerators requested: %d > %d", count, a.maxCountPerInstance,
		)
	}

	// Compute count as the next power of two
	shiftLeft := 63 - bits.LeadingZeros64(uint64(count))
	if bits.TrailingZeros64(uint64(count)) != shiftLeft {
		shiftLeft++
	}
	actualCount := 1 << shiftLeft

	// And return accelerator config
	return &computepb.AcceleratorConfig{
		AcceleratorType:  proto.String(a.uri),
		AcceleratorCount: proto.Int32(int32(actualCount)),
	}, nil
}

//-------------------------------------------------------------------------------------------------

// fetchAccelerators fetches all accelerators that are available in the provided zones. It then
// returns a list of available accelerators for each zone.
func fetchAccelerators(
	ctx context.Context, client *compute.AcceleratorTypesClient, project string, zones []string,
) (map[string][]Accelerator, error) {
	result := make(map[string][]Accelerator)
	for _, zone := range zones {
		result[zone] = make([]Accelerator, 0)
	}

	it := client.AggregatedList(ctx, &computepb.AggregatedListAcceleratorTypesRequest{
		Project: project,
	})
	err := gcputils.Iterate[compute.AcceleratorTypesScopedListPair](
		ctx, it, func(pair compute.AcceleratorTypesScopedListPair) error {
			zone := path.Base(pair.Key)
			if _, ok := result[zone]; !ok {
				// If the zone is not in the set of available zones, ignore
				return nil
			}

			for _, item := range pair.Value.AcceleratorTypes {
				// Ignore all virtual workstations
				if strings.HasSuffix(item.GetName(), "-vws") {
					continue
				}

				gpuKind, err := typedefs.GPUKindFromProviderGcp(item.GetName())
				if err != nil {
					// The GPU type could not be parsed, this might indicate an error
					zeus.Logger(ctx).Warn(
						"failed to parse GPU type", zap.String("zone", zone), zap.Error(err),
					)
					continue
				}

				// Otherwise we add the GPU to the permitted accelerators
				result[zone] = append(result[zone], Accelerator{
					uri:                 item.GetSelfLink(),
					kind:                gpuKind,
					maxCountPerInstance: uint16(item.GetMaximumCardsPerInstance()),
				})
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return result, nil
}
