//go:build integration

package gcpinstances

import "os"

var gcpProject = func() string {
	project := os.Getenv("GCP_PROJECT")
	if project == "" {
		panic("GCP_PROJECT environment variable is not set")
	}
	return project
}()
