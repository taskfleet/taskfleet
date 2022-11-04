package typedefs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// CPUArchitecture describes the architecture of a CPU.
type CPUArchitecture string

const (
	// ArchitectureX86 describes a CPU with x86-64 architecture.
	ArchitectureX86 CPUArchitecture = "x86-64"
	// ArchitectureArm describes a CPU with ARM64 architecture.
	ArchitectureArm CPUArchitecture = "arm64"
)

// UnmarshalJSON implements json.Unmarshaler.
func (a *CPUArchitecture) UnmarshalJSON(b []byte) error {
	var arch string
	if err := json.Unmarshal(b, &arch); err != nil {
		return err
	}
	switch arch {
	case string(ArchitectureX86):
		*a = ArchitectureX86
	case string(ArchitectureArm):
		*a = ArchitectureArm
	}
	return fmt.Errorf("invalid cpu architecture %q", arch)
}

//-------------------------------------------------------------------------------------------------

// ToProviderAws returns the AWS architecture type from the CPU architecture.
func (a CPUArchitecture) ToProviderAws() types.ArchitectureType {
	switch a {
	case ArchitectureArm:
		return types.ArchitectureTypeArm64
	case ArchitectureX86:
		return types.ArchitectureTypeX8664
	default:
		panic("invalid cpu architecture")
	}
}

// CPUArchitectureFromProviderAws returns the CPU architecture for the given AWS architecture.
func CPUArchitectureFromProviderAws(infos []types.ArchitectureType) (CPUArchitecture, error) {
	if len(infos) != 1 {
		return CPUArchitecture(""), fmt.Errorf("unable to handle multi-architecture CPU")
	}
	switch infos[0] {
	case types.ArchitectureTypeArm64:
		return ArchitectureArm, nil
	case types.ArchitectureTypeX8664:
		return ArchitectureX86, nil
	default:
		return CPUArchitecture(""), fmt.Errorf("unknown CPU architecture %q", infos[0])
	}
}
