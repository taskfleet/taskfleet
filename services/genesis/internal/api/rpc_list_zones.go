package api

import (
	"context"

	genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
)

// ListZones implements the genesis interface.
func (s *Service) ListZones(
	ctx context.Context, request *genesis_v1.ListZonesRequest,
) (*genesis_v1.ListZonesResponse, error) {
	result := make([]*genesis_v1.Zone, 0)
	for _, zone := range s.provider.Zones().List() {
		gpuKinds := make([]genesis_v1.GPUKind, len(zone.GPUs))
		for i := range zone.GPUs {
			gpuKinds[i] = zone.GPUs[i].MarshalProto()
		}

		result = append(result, &genesis_v1.Zone{
			Provider:      s.provider.CloudProvider().MarshalProto(),
			Name:          zone.Name,
			AvailableGpus: gpuKinds,
		})
	}
	return &genesis_v1.ListZonesResponse{Zones: result}, nil
}
