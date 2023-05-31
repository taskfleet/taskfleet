package api

import (
	"fmt"

	"github.com/docker/go-units"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

type diskSizes struct {
	hdd                uint32
	ssdStandard        uint32
	ssdHighPerformance uint32
}

func parseDisks(
	templates []v1alpha1.InstanceDisk, cpuCount uint32,
) ([]providers.DiskConfig, diskSizes, error) {
	disks := []providers.DiskConfig{}
	sizes := diskSizes{}

	for i, template := range templates {
		// Get size
		sizePerCPU, err := units.RAMInBytes(string(template.SizePerCPU))
		if err != nil {
			return nil, diskSizes{},
				fmt.Errorf("Failed to parse size per CPU for template %d: %s", i, err)
		}
		size := uint32(sizePerCPU/(1024*1024*1024)) * cpuCount

		// Create disk
		switch template.Type {
		case v1alpha1.DiskTypeHDD:
			sizes.hdd += size
			disks = append(disks, providers.DiskConfig{
				Name: template.Name, Type: typedefs.DiskHDD, SizeGiB: size,
			})
		case v1alpha1.DiskTypeSSDStandard:
			sizes.ssdStandard += size
			disks = append(disks, providers.DiskConfig{
				Name: template.Name, Type: typedefs.DiskSSDStandard, SizeGiB: size,
			})
		case v1alpha1.DiskTypeSSDHighPerformance:
			sizes.ssdHighPerformance += size
			disks = append(disks, providers.DiskConfig{
				Name: template.Name, Type: typedefs.DiskSSDHighPerformance, SizeGiB: size,
			})
		}
	}
	return disks, sizes, nil
}
