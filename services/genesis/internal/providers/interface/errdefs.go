package providers

import (
	"errors"
	"net/http"

	"google.golang.org/api/googleapi"
)

// IsErrNotFound returns whether the given error describes an API error which indicates that a
// resource could not be found.
func IsErrNotFound(err error) bool {
	var googleApiError *googleapi.Error
	if errors.As(err, &googleApiError) {
		return googleApiError.Code == http.StatusNotFound
	}
	return false
}
