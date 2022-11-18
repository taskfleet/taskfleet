package instances

func validateMemoryPerCPU(cpus uint16, memoryMib uint32) bool {
	mem := float64(memoryMib) / 1024
	memPerCPU := mem / float64(cpus)

	// For more than 640 GiB of RAM, we allow between 12 and 24 GiB of RAM per CPU.
	if mem > 640 {
		return memPerCPU >= 12 && memPerCPU <= 24
	}

	// For more than 320 GiB of RAM, we allow between 6 and 12 GiB of RAM per CPU.
	if mem > 320 {
		return memPerCPU >= 6 && memPerCPU <= 12
	}

	// For more than 80 GiB of RAM, we allow between 2 and 10 GiB of RAM per CPU.
	if mem > 80 {
		return memPerCPU >= 2 && memPerCPU <= 10
	}

	// Otherwise, we allow between 1 and 8 GiB per CPU.
	return memPerCPU >= 1 && memPerCPU <= 8
}
