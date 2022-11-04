package providers

import "fmt"

// // InstanceRefFromCommonName returns an instance ref from the given common name and a zone.
// func InstanceRefFromCommonName(name string, zone string) (InstanceMeta, error) {
// 	id, err := uuid.Parse(strings.TrimPrefix(name, "taskfleet-"))
// 	if err != nil {
// 		return InstanceMeta{}, err
// 	}
// 	return InstanceMeta{ID: id, Zone: zone}, nil
// }

// CommonName returns a common name for the provided instance reference.
func (r InstanceMeta) CommonName() string {
	return fmt.Sprintf("taskfleet-%s", r.ID)
}
