package gcpinstances

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestNewDisksHelper(t *testing.T) {
	ctx := context.Background()
	linkMap := map[string]string{
		"europe-west3-c": "https://example.com/europe-west3-c/pd-balanced",
		"us-east1-c":     "https://example.com/us-east1-c/pd-balanced",
	}
	client := newDiskTypesClient(ctx, t, linkMap)
	helper, err := newDisksHelper(
		ctx,
		"",
		template.GcpBootConfig{DiskSize: "10Gi"},
		[]template.InstanceDisk{{Name: "disk1", SizePerCPU: "5Gi"}},
		"pd-balanced",
		client,
	)
	assert.Nil(t, err)

	// Assert assumptions
	assert.EqualValues(t, helper.bootDiskSizeGiB, 10)
	assert.ElementsMatch(t, helper.extraDisks, []disk{{name: "disk1", sizePerCpuGiB: 5}})
	assert.True(t, reflect.DeepEqual(helper.diskTypeSelfLinks, linkMap))
}

func TestDisksHelperDiskConfig(t *testing.T) {
	linkMap := map[string]string{
		"europe-west3-c": "https://example.com/europe-west3-c/pd-balanced",
		"us-east1-c":     "https://example.com/us-east1-c/pd-balanced",
	}
	helper := &disksHelper{
		bootDiskSizeGiB:   10,
		bootImages:        []template.Option[string]{{Config: "my-boot-image"}},
		extraDisks:        []disk{{sizePerCpuGiB: 5}},
		diskTypeSelfLinks: linkMap,
	}

	// Run disk config
	disks := helper.diskConfig(
		"my-instance",
		"europe-west3-c",
		instances.Resources{CPUCount: 3},
		typedefs.ArchitectureX86,
	)
	require.Len(t, disks, 2)

	bootDisk := disks[0]
	assert.True(t, bootDisk.GetAutoDelete())
	assert.True(t, bootDisk.GetBoot())
	assert.EqualValues(t, bootDisk.InitializeParams.GetDiskSizeGb(), 10)
	assert.Equal(t, bootDisk.InitializeParams.GetDiskName(), "my-instance-boot-disk")
	assert.Equal(
		t,
		bootDisk.InitializeParams.GetDiskType(),
		"https://example.com/europe-west3-c/pd-balanced",
	)
	assert.Equal(t, bootDisk.InitializeParams.GetSourceImage(), "my-boot-image")

	extraDisk := disks[1]
	assert.True(t, extraDisk.GetAutoDelete())
	assert.False(t, extraDisk.GetBoot())
	assert.EqualValues(t, extraDisk.InitializeParams.GetDiskSizeGb(), 15)
	assert.Equal(
		t, extraDisk.InitializeParams.GetDiskType(),
		"https://example.com/europe-west3-c/pd-balanced",
	)
	assert.Equal(t, extraDisk.InitializeParams.GetSourceImage(), "")
}

func TestNewReservationsHelper(t *testing.T) {
	testCases := []struct {
		value    *string
		expected uint32
		isErr    bool
	}{
		{value: jack.Ptr("100Mi"), expected: 100},
		{value: jack.Ptr("4Gi"), expected: 4096},
		{value: nil, expected: 0},
		{value: jack.Ptr("1X"), isErr: true},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			helper, err := newReservationsHelper(template.InstanceReservations{
				Memory: tc.value,
			})
			assert.Equal(t, tc.isErr, err != nil)
			if err == nil {
				fmt.Println(tc.expected, helper.memoryMiB)
				assert.Equal(t, tc.expected, helper.memoryMiB)
			}
		})
	}
}

func TestReservationsHelperUpdateResources(t *testing.T) {
	helper := &reservationsHelper{memoryMiB: 256}
	r := helper.updateResources(instances.Resources{MemoryMiB: 1024})
	assert.EqualValues(t, r.MemoryMiB, 1280)
}

//-------------------------------------------------------------------------------------------------

func newDiskTypesClient(
	ctx context.Context, t *testing.T, diskLinks map[string]string,
) *compute.DiskTypesClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if diskLinks == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response := &computepb.DiskTypeAggregatedList{
			Items: map[string]*computepb.DiskTypesScopedList{},
		}
		for zone, diskLink := range diskLinks {
			response.Items[zone] = &computepb.DiskTypesScopedList{DiskTypes: []*computepb.DiskType{
				{SelfLink: proto.String(diskLink)},
			}}
		}
		result := jack.Must(protojson.Marshal(response))
		jack.Must(w.Write(result))
	}))
	service, err := compute.NewDiskTypesRESTClient(
		ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication(),
	)
	require.Nil(t, err)
	return service
}
