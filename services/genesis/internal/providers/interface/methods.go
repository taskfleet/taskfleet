package providers

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// InstanceRefFromCommonName returns an instance ref from the given common name and a zone.
func InstanceMetaFromCommonName(name string, zone string) (InstanceMeta, error) {
	id, err := uuid.Parse(strings.TrimPrefix(name, "taskfleet-"))
	if err != nil {
		return InstanceMeta{}, err
	}
	return InstanceMeta{ID: id, ProviderZone: zone}, nil
}

// CommonName returns a common name for the provided instance reference.
func (r InstanceMeta) CommonName() string {
	return fmt.Sprintf("taskfleet-%s", r.ID)
}
