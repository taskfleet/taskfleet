// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: genesis/messages/v1/instance_event.proto

package genesis_messages

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// Validate checks the field values on InstanceEvent with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *InstanceEvent) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetInstance()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InstanceEventValidationError{
				field:  "Instance",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetTimestamp()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InstanceEventValidationError{
				field:  "Timestamp",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	switch m.Event.(type) {

	case *InstanceEvent_Created:

		if v, ok := interface{}(m.GetCreated()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InstanceEventValidationError{
					field:  "Created",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *InstanceEvent_CreationFailed:

		if v, ok := interface{}(m.GetCreationFailed()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InstanceEventValidationError{
					field:  "CreationFailed",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *InstanceEvent_Deleted:

		if v, ok := interface{}(m.GetDeleted()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return InstanceEventValidationError{
					field:  "Deleted",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// InstanceEventValidationError is the validation error returned by
// InstanceEvent.Validate if the designated constraints aren't met.
type InstanceEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InstanceEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InstanceEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InstanceEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InstanceEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InstanceEventValidationError) ErrorName() string { return "InstanceEventValidationError" }

// Error satisfies the builtin error interface
func (e InstanceEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInstanceEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InstanceEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InstanceEventValidationError{}

// Validate checks the field values on InstanceCreatedEvent with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *InstanceCreatedEvent) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InstanceCreatedEventValidationError{
				field:  "Config",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetResources()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return InstanceCreatedEventValidationError{
				field:  "Resources",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Hostname

	return nil
}

// InstanceCreatedEventValidationError is the validation error returned by
// InstanceCreatedEvent.Validate if the designated constraints aren't met.
type InstanceCreatedEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InstanceCreatedEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InstanceCreatedEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InstanceCreatedEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InstanceCreatedEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InstanceCreatedEventValidationError) ErrorName() string {
	return "InstanceCreatedEventValidationError"
}

// Error satisfies the builtin error interface
func (e InstanceCreatedEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInstanceCreatedEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InstanceCreatedEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InstanceCreatedEventValidationError{}

// Validate checks the field values on InstanceCreationFailedEvent with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *InstanceCreationFailedEvent) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Reason

	// no validation rules for Message

	return nil
}

// InstanceCreationFailedEventValidationError is the validation error returned
// by InstanceCreationFailedEvent.Validate if the designated constraints
// aren't met.
type InstanceCreationFailedEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InstanceCreationFailedEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InstanceCreationFailedEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InstanceCreationFailedEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InstanceCreationFailedEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InstanceCreationFailedEventValidationError) ErrorName() string {
	return "InstanceCreationFailedEventValidationError"
}

// Error satisfies the builtin error interface
func (e InstanceCreationFailedEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInstanceCreationFailedEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InstanceCreationFailedEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InstanceCreationFailedEventValidationError{}

// Validate checks the field values on InstanceDeletedEvent with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *InstanceDeletedEvent) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Reason

	return nil
}

// InstanceDeletedEventValidationError is the validation error returned by
// InstanceDeletedEvent.Validate if the designated constraints aren't met.
type InstanceDeletedEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e InstanceDeletedEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e InstanceDeletedEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e InstanceDeletedEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e InstanceDeletedEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e InstanceDeletedEventValidationError) ErrorName() string {
	return "InstanceDeletedEventValidationError"
}

// Error satisfies the builtin error interface
func (e InstanceDeletedEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sInstanceDeletedEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = InstanceDeletedEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = InstanceDeletedEventValidationError{}
