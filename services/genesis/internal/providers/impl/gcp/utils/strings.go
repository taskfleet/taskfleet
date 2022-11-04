package gcputils

import (
	"strings"
)

// RegionFromZone returns the region in which the provided zone is located.
func RegionFromZone(zone string) string {
	return strings.Join(strings.Split(zone, "-")[:2], "-")
}
