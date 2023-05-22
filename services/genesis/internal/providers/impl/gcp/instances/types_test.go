package gcpinstances

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestFindAvailableInstanceTypes(t *testing.T) {
	ctx := context.Background()
	client := newMachineTypesClient(ctx, t, map[string][]*computepb.MachineType{
		"europe-west3-b": {
			{Name: proto.String("n1-standard-1"), GuestCpus: proto.Int32(1)},
			{
				Name:      proto.String("a2-highgpu-1"),
				GuestCpus: proto.Int32(8),
				Accelerators: []*computepb.Accelerators{
					{
						GuestAcceleratorType:  proto.String("nvidia-tesla-a100"),
						GuestAcceleratorCount: proto.Int32(1),
					},
				},
			},
		},
		"us-east1-c": {
			{Name: proto.String("n1-standard-2"), GuestCpus: proto.Int32(2)},
		},
		"unknown": {
			{Name: proto.String("n1-standard-4"), GuestCpus: proto.Int32(4)},
		},
	})
	zones := gcpzones.NewMockClient(t)
	zones.EXPECT().List().Return([]providers.Zone{
		{Name: "europe-west3-b", GPUs: []typedefs.GPUKind{typedefs.GPUNvidiaTeslaA100}},
		{Name: "us-east1-c", GPUs: []typedefs.GPUKind{typedefs.GPUNvidiaTeslaK80}},
	})
	zones.EXPECT().GetAccelerator("europe-west3-b", typedefs.GPUNvidiaTeslaA100).Return(
		gcpzones.Accelerator{Kind: typedefs.GPUNvidiaTeslaA100, MaxCountPerInstance: 1}, nil,
	)
	zones.EXPECT().GetAccelerator("us-east1-c", typedefs.GPUNvidiaTeslaK80).Return(
		gcpzones.Accelerator{Kind: typedefs.GPUNvidiaTeslaK80, MaxCountPerInstance: 2}, nil,
	)

	types, err := findAvailableInstanceTypes(ctx, client, zones, "")
	assert.Nil(t, err)

	assert.ElementsMatch(t, jack.MapKeys(types), []string{"europe-west3-b", "us-east1-c"})
	assert.ElementsMatch(t, jack.SliceMap(types["europe-west3-b"], func(t instances.Type) string {
		return t.Name
	}), []string{"n1-standard-1", "a2-highgpu-1"})
	assert.ElementsMatch(t, jack.SliceMap(types["us-east1-c"], func(t instances.Type) string {
		return t.Name
	}), []string{
		"n1-standard-2", "n1-standard-2-nvidia-tesla-k80-1", "n1-standard-2-nvidia-tesla-k80-2",
	})
}

func TestTryUnmarshalInstanceType(t *testing.T) {
	testCases := []struct {
		input    *computepb.MachineType
		expected *instances.Type
	}{
		{ // No shared CPU
			input:    &computepb.MachineType{IsSharedCpu: proto.Bool(true)},
			expected: nil,
		},
		{ // No deprecated instances
			input:    &computepb.MachineType{Deprecated: &computepb.DeprecationStatus{}},
			expected: nil,
		},
		{ // No m2 instances
			input:    &computepb.MachineType{Name: proto.String("m2-highmem-16")},
			expected: nil,
		},
		{
			// X86 instances
			input: &computepb.MachineType{
				Name:      proto.String("n1-standard-1"),
				SelfLink:  proto.String("https://example.com/n1-standard-1"),
				GuestCpus: proto.Int32(1),
				MemoryMb:  proto.Int32(3500),
			},
			expected: &instances.Type{
				Name:         "n1-standard-1",
				UID:          "https://example.com/n1-standard-1",
				Architecture: typedefs.ArchitectureX86,
				Resources: instances.Resources{
					CPUCount:  1,
					MemoryMiB: 3500,
				},
			},
		},
		{
			// ARM instance
			input: &computepb.MachineType{
				Name:      proto.String("t2a-standard-1"),
				SelfLink:  proto.String("https://example.com/t2a-standard-1"),
				GuestCpus: proto.Int32(1),
				MemoryMb:  proto.Int32(4096),
			},
			expected: &instances.Type{
				Name:         "t2a-standard-1",
				UID:          "https://example.com/t2a-standard-1",
				Architecture: typedefs.ArchitectureArm,
				Resources: instances.Resources{
					CPUCount:  1,
					MemoryMiB: 4096,
				},
			},
		},
		{
			// GPU instance
			input: &computepb.MachineType{
				Name:      proto.String("a2-highgpu-1"),
				SelfLink:  proto.String("https://example.com/a2-highgpu-1"),
				GuestCpus: proto.Int32(8),
				MemoryMb:  proto.Int32(16384),
				Accelerators: []*computepb.Accelerators{
					{
						GuestAcceleratorCount: proto.Int32(1),
						GuestAcceleratorType:  proto.String("nvidia-tesla-a100"),
					},
				},
			},
			expected: &instances.Type{
				Name:         "a2-highgpu-1",
				UID:          "https://example.com/a2-highgpu-1",
				Architecture: typedefs.ArchitectureX86,
				Resources: instances.Resources{
					CPUCount:  8,
					MemoryMiB: 16384,
					GPU: &instances.GPUResources{
						Kind:  typedefs.GPUNvidiaTeslaA100,
						Count: 1,
					},
				},
			},
		},
		{
			// Unknown machine type family
			input:    &computepb.MachineType{Name: proto.String("xxx-unknown-1")},
			expected: nil,
		},
		{
			// Unparseable accelerator
			input: &computepb.MachineType{
				Name: proto.String("a2-unknown-1"),
				Accelerators: []*computepb.Accelerators{
					{
						GuestAcceleratorCount: proto.Int32(1),
						GuestAcceleratorType:  proto.String("xxx"),
					},
				},
			},
			expected: nil,
		},
	}

	ctx := context.Background()
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := tryUnmarshalInstanceType(ctx, tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestExplodeAvailableGpuInstanceTypes(t *testing.T) {
	n1Instances := []instances.Type{
		{
			Name:         "n1-standard-1",
			UID:          "https://example.com/n1-standard-1",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  1,
				MemoryMiB: 3500,
			},
		},
		{
			Name:         "n1-standard-4",
			UID:          "https://example.com/n1-standard-4",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  4,
				MemoryMiB: 14000,
			},
		},
	}

	actual := explodeAvailableGpuInstanceTypes(
		"my-zone", n1Instances, typedefs.GPUNvidiaTeslaK80, 4,
	)

	expected := []instances.Type{
		{
			Name:         "n1-standard-1-nvidia-tesla-k80-1",
			UID:          "https://example.com/n1-standard-1",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  1,
				MemoryMiB: 3500,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 1},
			},
		},
		{
			Name:         "n1-standard-1-nvidia-tesla-k80-2",
			UID:          "https://example.com/n1-standard-1",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  1,
				MemoryMiB: 3500,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 2},
			},
		},
		{
			Name:         "n1-standard-1-nvidia-tesla-k80-4",
			UID:          "https://example.com/n1-standard-1",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  1,
				MemoryMiB: 3500,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 4},
			},
		},
		{
			Name:         "n1-standard-4-nvidia-tesla-k80-1",
			UID:          "https://example.com/n1-standard-4",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  4,
				MemoryMiB: 14000,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 1},
			},
		},
		{
			Name:         "n1-standard-4-nvidia-tesla-k80-2",
			UID:          "https://example.com/n1-standard-4",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  4,
				MemoryMiB: 14000,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 2},
			},
		},
		{
			Name:         "n1-standard-4-nvidia-tesla-k80-4",
			UID:          "https://example.com/n1-standard-4",
			Architecture: typedefs.ArchitectureX86,
			Resources: instances.Resources{
				CPUCount:  4,
				MemoryMiB: 14000,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 4},
			},
		},
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestExtendedInstanceTypeName(t *testing.T) {
	testCases := []struct {
		instance string
		gpu      *instances.GPUResources
		expected string
	}{
		{instance: "n1-standard-1", gpu: nil, expected: "n1-standard-1"},
		{instance: "n2-standard-1", gpu: nil, expected: "n2-standard-1"},
		{
			instance: "a2-highgpu-1",
			gpu:      &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaA10, Count: 1},
			expected: "a2-highgpu-1",
		},
		{
			instance: "n1-standard-1",
			gpu:      &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 1},
			expected: "n1-standard-1-nvidia-tesla-k80-1",
		},
		{
			instance: "n1-standard-4",
			gpu:      &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaV100, Count: 2},
			expected: "n1-standard-4-nvidia-tesla-v100-2",
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			actual := extendedInstanceTypeName(tc.instance, tc.gpu)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestMaxCpuAndMemoryForGpu(t *testing.T) {
	testCases := []struct {
		gpuKind      typedefs.GPUKind
		gpuCount     uint16
		zone         string
		maxCpuCount  uint16
		maxMemoryGiB uint32
	}{
		{gpuKind: typedefs.GPUNvidiaTeslaT4, gpuCount: 1, maxCpuCount: 48, maxMemoryGiB: 312},
		{gpuKind: typedefs.GPUNvidiaTeslaT4, gpuCount: 2, maxCpuCount: 48, maxMemoryGiB: 312},
		{gpuKind: typedefs.GPUNvidiaTeslaT4, gpuCount: 4, maxCpuCount: 96, maxMemoryGiB: 624},
		{gpuKind: typedefs.GPUNvidiaTeslaP4, gpuCount: 1, maxCpuCount: 24, maxMemoryGiB: 156},
		{gpuKind: typedefs.GPUNvidiaTeslaP4, gpuCount: 2, maxCpuCount: 48, maxMemoryGiB: 312},
		{gpuKind: typedefs.GPUNvidiaTeslaV100, gpuCount: 1, maxCpuCount: 12, maxMemoryGiB: 78},
		{gpuKind: typedefs.GPUNvidiaTeslaV100, gpuCount: 2, maxCpuCount: 24, maxMemoryGiB: 156},
		{gpuKind: typedefs.GPUNvidiaTeslaP100, gpuCount: 1, maxCpuCount: 16, maxMemoryGiB: 104},
		{gpuKind: typedefs.GPUNvidiaTeslaP100, gpuCount: 2, maxCpuCount: 32, maxMemoryGiB: 208},
		{
			gpuKind:      typedefs.GPUNvidiaTeslaP100,
			gpuCount:     4,
			maxCpuCount:  64,
			maxMemoryGiB: 208,
			zone:         "us-east1-c",
		},
		{
			gpuKind:      typedefs.GPUNvidiaTeslaP100,
			gpuCount:     4,
			maxCpuCount:  96,
			maxMemoryGiB: 624,
			zone:         "us-west1-c",
		},
		{gpuKind: typedefs.GPUNvidiaTeslaK80, gpuCount: 1, maxCpuCount: 8, maxMemoryGiB: 52},
		{gpuKind: typedefs.GPUNvidiaTeslaK80, gpuCount: 2, maxCpuCount: 16, maxMemoryGiB: 104},
		{gpuKind: typedefs.GPUNvidiaTeslaK80, gpuCount: 4, maxCpuCount: 32, maxMemoryGiB: 208},
		{
			gpuKind:      typedefs.GPUNvidiaTeslaK80,
			gpuCount:     8,
			maxCpuCount:  64,
			maxMemoryGiB: 416,
			zone:         "asia-east1-a",
		},
		{
			gpuKind:      typedefs.GPUNvidiaTeslaK80,
			gpuCount:     8,
			maxCpuCount:  64,
			maxMemoryGiB: 208,
			zone:         "us-west1-c",
		},
		{gpuKind: typedefs.GPUNvidiaTeslaA100Gb80, gpuCount: 1, maxCpuCount: 0, maxMemoryGiB: 0},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cpu, mem := maxCpuAndMemoryForGpu(tc.gpuKind, tc.gpuCount, tc.zone)
			assert.Equalf(
				t, tc.maxCpuCount, cpu,
				"max CPU does not match for %d %s", tc.gpuCount, tc.gpuKind,
			)
			assert.Equal(
				t, tc.maxMemoryGiB, mem,
				"max memory does not match for %d %s", tc.gpuCount, tc.gpuKind,
			)
		})

	}
}

//-------------------------------------------------------------------------------------------------

func newMachineTypesClient(
	ctx context.Context, t *testing.T, machineTypes map[string][]*computepb.MachineType,
) *compute.MachineTypesClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if machineTypes == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response := &computepb.MachineTypeAggregatedList{
			Items: map[string]*computepb.MachineTypesScopedList{},
		}
		for zone, machines := range machineTypes {
			response.Items[zone] = &computepb.MachineTypesScopedList{MachineTypes: machines}
		}
		result := jack.Must(protojson.Marshal(response))
		jack.Must(w.Write(result))
	}))
	service, err := compute.NewMachineTypesRESTClient(
		ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication(),
	)
	require.Nil(t, err)
	return service
}
