package awsutils

// RegionFromZone returns the region in which the provided availability zone is located.
func RegionFromZone(zone string) string {
	return zone[:len(zone)-1]
}
