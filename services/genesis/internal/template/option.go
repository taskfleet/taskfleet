package template

import (
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// Option is a type that simplifies iterating over multiple options, selecting one whose selector
// matches a query.
type Option[T any] struct {
	// The actual underlying configuration for the option.
	Config T `json:"config"`
	// The selector to use that ought to match the query.
	Selector OptionSelector `json:"selector,omitempty"`
}

// OptionSelector defines desired values of instance properties.
type OptionSelector struct {
	// A set of GPUs that the instance template matches. If not provided, matches all GPU kinds and
	GPUs *Selector[*typedefs.GPUKind] `json:"gpus,omitempty"`
	// A set of CPU architectures that the instance template matches. If not provided, matches all
	// CPU architectures.
	CPUs *Selector[typedefs.CPUArchitecture] `json:"cpus,omitempty"`
}

// Selector defines the matching expressions for limiting configurations to a particular resource.
type Selector[T any] struct {
	// The resources for which the template applies.
	In []T `json:"in,omitempty"`
}

//-------------------------------------------------------------------------------------------------

type selectorType interface {
	*typedefs.GPUKind | typedefs.CPUArchitecture
}

// MatchingOption selects from the given options the first one whose selector matches the provided
// GPU kind and CPU architecture. If no such option is found, `nil` will be returned. Otherwise,
// the option's underlying configuration will be returned.
func MatchingOption[T any](
	options []Option[T], gpu *typedefs.GPUKind, cpu typedefs.CPUArchitecture,
) *T {
	for _, option := range options {
		if option.Selector.GPUs != nil {
			if !contains(option.Selector.GPUs.In, gpu) {
				continue
			}
		}
		if option.Selector.CPUs != nil {
			if !contains(option.Selector.CPUs.In, cpu) {
				continue
			}
		}
		return &option.Config
	}
	return nil
}

func contains[T selectorType](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
