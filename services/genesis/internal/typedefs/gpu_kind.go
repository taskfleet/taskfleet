package typedefs

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
)

// GPUKind describes a GPU type.
type GPUKind string

const (
	// GPUNvidiaTeslaK80 represents a Tesla K80 GPU.
	GPUNvidiaTeslaK80 GPUKind = "nvidia-tesla-k80"
	// GPUNvidiaTeslaM60 represents a Tesla M60 GPU.
	GPUNvidiaTeslaM60 GPUKind = "nvidia-tesla-m60"
	// GPUNvidiaTeslaP100 represents a Tesla P100 GPU.
	GPUNvidiaTeslaP100 GPUKind = "nvidia-tesla-p100"
	// GPUNvidiaTeslaP4 represents a Tesla P4 GPU.
	GPUNvidiaTeslaP4 GPUKind = "nvidia-tesla-p4"
	// GPUNvidiaTeslaV100 represents a Tesla V100 GPU.
	GPUNvidiaTeslaV100 GPUKind = "nvidia-tesla-v100"
	// GPUNvidiaTeslaT4 represents a Tesla T4 GPU.
	GPUNvidiaTeslaT4 GPUKind = "nvidia-tesla-t4"
	// GPUNvidiaTeslaA100 represents a Tesla A100 GPU with 40 GB of memory.
	GPUNvidiaTeslaA100 GPUKind = "nvidia-tesla-a100"
	// GPUNvidiaTeslaA10 represents a Tesla A10 GPU.
	GPUNvidiaTeslaA10 GPUKind = "nvidia-tesla-a10"
	// GPUNvidiaTeslaA100Gb80 represents a Tesla A100 GPU with 80 GB of memory.
	GPUNvidiaTeslaA100Gb80 GPUKind = "nvidia-tesla-a100-80gb"
)

// UnmarshalJSON implements json.Unmarshaler.
func (k *GPUKind) UnmarshalJSON(b []byte) error {
	var kind string
	if err := json.Unmarshal(b, &kind); err != nil {
		return err
	}
	switch kind {
	case string(GPUNvidiaTeslaK80):
		*k = GPUNvidiaTeslaK80
	case string(GPUNvidiaTeslaM60):
		*k = GPUNvidiaTeslaM60
	case string(GPUNvidiaTeslaP100):
		*k = GPUNvidiaTeslaP100
	case string(GPUNvidiaTeslaP4):
		*k = GPUNvidiaTeslaP4
	case string(GPUNvidiaTeslaV100):
		*k = GPUNvidiaTeslaV100
	case string(GPUNvidiaTeslaT4):
		*k = GPUNvidiaTeslaT4
	case string(GPUNvidiaTeslaA100):
		*k = GPUNvidiaTeslaA100
	case string(GPUNvidiaTeslaA10):
		*k = GPUNvidiaTeslaA10
	case string(GPUNvidiaTeslaA100Gb80):
		*k = GPUNvidiaTeslaA100Gb80
	}
	return fmt.Errorf("invalid gpu kind %q", kind)
}

//-------------------------------------------------------------------------------------------------

// GPUKindUnmarshalProto returns the internal GPU representation for the provided hermes GPU kind.
func GPUKindUnmarshalProto(message genesis_v1.GPUKind) GPUKind {
	switch message {
	case genesis_v1.GPUKind_GPU_KIND_TESLA_K80:
		return GPUNvidiaTeslaK80
	case genesis_v1.GPUKind_GPU_KIND_TESLA_M60:
		return GPUNvidiaTeslaM60
	case genesis_v1.GPUKind_GPU_KIND_TESLA_P100:
		return GPUNvidiaTeslaP100
	case genesis_v1.GPUKind_GPU_KIND_TESLA_P4:
		return GPUNvidiaTeslaP4
	case genesis_v1.GPUKind_GPU_KIND_TESLA_V100:
		return GPUNvidiaTeslaV100
	case genesis_v1.GPUKind_GPU_KIND_TESLA_T4:
		return GPUNvidiaTeslaT4
	case genesis_v1.GPUKind_GPU_KIND_TESLA_A100:
		return GPUNvidiaTeslaA100
	case genesis_v1.GPUKind_GPU_KIND_TESLA_A10:
		return GPUNvidiaTeslaA10
	case genesis_v1.GPUKind_GPU_KIND_TESLA_A100_80GB:
		return GPUNvidiaTeslaA100Gb80
	default:
		panic("unknown GPU kind")
	}
}

// MarshalProto returns the hermes enum of the GPU kind.
func (k GPUKind) MarshalProto() genesis_v1.GPUKind {
	switch k {
	case GPUNvidiaTeslaK80:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_K80
	case GPUNvidiaTeslaM60:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_M60
	case GPUNvidiaTeslaP100:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_P100
	case GPUNvidiaTeslaP4:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_P4
	case GPUNvidiaTeslaV100:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_V100
	case GPUNvidiaTeslaT4:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_T4
	case GPUNvidiaTeslaA100:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_A100
	case GPUNvidiaTeslaA10:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_A10
	case GPUNvidiaTeslaA100Gb80:
		return genesis_v1.GPUKind_GPU_KIND_TESLA_A100_80GB
	default:
		panic("unknown GPU kind")
	}
}

// ShortName returns the short name of the GPU type.
func (k GPUKind) ShortName() string {
	splits := strings.Split(string(k), "-")
	return splits[len(splits)-1]
}

//-------------------------------------------------------------------------------------------------

// GPUKindFromProviderAws returns the GPU kind for the specified AWS GPU kind and an error if the
// AWS GPU is unknown.
func GPUKindFromProviderAws(info types.GpuDeviceInfo) (GPUKind, error) {
	if info.Manufacturer == nil || info.Name == nil {
		return GPUKind(""), fmt.Errorf("gpu info misses at least one of manufacturer or name")
	}
	if *info.Manufacturer != "NVIDIA" {
		return GPUKind(""), fmt.Errorf(
			"gpu manufacturer is other than NVIDIA is not supported (found %q)",
			*info.Manufacturer,
		)
	}
	switch *info.Name {
	case "K80":
		return GPUNvidiaTeslaK80, nil
	case "M60":
		return GPUNvidiaTeslaM60, nil
	case "T4", "T4g":
		return GPUNvidiaTeslaT4, nil
	case "V100":
		return GPUNvidiaTeslaV100, nil
	case "A10G":
		return GPUNvidiaTeslaA10, nil
	case "A100":
		return GPUNvidiaTeslaA100, nil
	default:
		return GPUKind(""), fmt.Errorf(
			"gpu info references unknown NVIDIA GPU (name %q)", *info.Name,
		)
	}
}

// GPUKindFromProviderGcp returns the GPU kind for the specified GCP GPU kind and an error if the
// GCP GPU is not available.
func GPUKindFromProviderGcp(kind string) (GPUKind, error) {
	switch kind {
	case "nvidia-tesla-k80":
		return GPUNvidiaTeslaK80, nil
	case "nvidia-tesla-p4":
		return GPUNvidiaTeslaP4, nil
	case "nvidia-tesla-p100":
		return GPUNvidiaTeslaP100, nil
	case "nvidia-tesla-t4":
		return GPUNvidiaTeslaT4, nil
	case "nvidia-tesla-v100":
		return GPUNvidiaTeslaV100, nil
	case "nvidia-tesla-a100":
		return GPUNvidiaTeslaA100, nil
	case "nvidia-a100-80gb":
		return GPUNvidiaTeslaA100Gb80, nil
	default:
		return GPUKind(""), fmt.Errorf("gpu kind %q is unknown for GCP", kind)
	}
}
