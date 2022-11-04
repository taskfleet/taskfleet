package providers

import (
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi"
)

// IsErrNotFound returns whether the given error describes an API error which indicates that a
// resource could not be found.
func IsErrNotFound(err error) bool {
	switch err := err.(type) {
	case APIError:
		switch err := err.Cause.(type) {
		case *googleapi.Error:
			if err.Code == http.StatusNotFound {
				return true
			}
		}
	}
	return false
}

//-------------------------------------------------------------------------------------------------
// CLIENT ERROR
//-------------------------------------------------------------------------------------------------

// ClientError is an error that is returned whenever retrying a call will have no effect since the
// parameters specified by the client are invalid.
type ClientError struct {
	Message string
	Cause   error
}

// NewClientError initializes a new client error.
func NewClientError(message string, cause error) ClientError {
	return ClientError{Message: message, Cause: cause}
}

func (e ClientError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("client error (%s): %s", e.Message, e.Cause)
	}
	return fmt.Sprintf("client error (%s)", e.Message)
}

func (e ClientError) String() string {
	return e.Error()
}

//-------------------------------------------------------------------------------------------------
// API ERROR
//-------------------------------------------------------------------------------------------------

// APIError is an error that is returned whenever an API call that is expected to succeed fails.
// Retrying could help, however, a long enough period should be waited.
type APIError struct {
	Message string
	Cause   error
}

// NewAPIError initializes a new API error.
func NewAPIError(message string, cause error) APIError {
	return APIError{Message: message, Cause: cause}
}

func (e APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("api error (%s): %s", e.Message, e.Cause)
	}
	return fmt.Sprintf("api error (%s)", e.Message)
}

func (e APIError) String() string {
	return e.Error()
}

//-------------------------------------------------------------------------------------------------
// FATAL ERROR
//-------------------------------------------------------------------------------------------------

// FatalError is an error that is returned whenever some internal invariant is violated. This
// usually means that something is quite wrong.
type FatalError struct {
	Message string
	Cause   error
}

// NewFatalError initializes a new fatal error.
func NewFatalError(message string, cause error) FatalError {
	return FatalError{Message: message, Cause: cause}
}

func (e FatalError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("fatal error (%s): %s", e.Message, e.Cause)
	}
	return fmt.Sprintf("fatal error (%s)", e.Message)
}

func (e FatalError) String() string {
	return e.Error()
}
