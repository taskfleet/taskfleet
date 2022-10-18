package eagle

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	errInvalidType = errors.New("invalid type encountered")
)

// IsErrInvalidType checks whether the provided error indicates that the encountered type was
// invalid when unmarshalling a string.
func IsErrInvalidType(err error) bool {
	return errors.Is(err, errInvalidType)
}

//-------------------------------------------------------------------------------------------------

// String is a type which can be used in configurations. When unmarshalling a JSON configuration,
// it either accepts a value `"string"` or an object `{"file": "filename"}`. When a file is given,
// unmarshalling reads the file and stores its contents.
type String struct {
	value string
}

// NewString creates a new string using the specified value.
func NewString(value string) String {
	return String{value}
}

// NewStringFromFile creates a new string by reading the specified file and fails if reading fails.
func NewStringFromFile(file string) (String, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return String{}, err
	}
	return String{string(data)}, nil
}

// Value returns the value of the string.
func (s String) Value() string {
	return s.value
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *String) UnmarshalJSON(b []byte) error {
	var element interface{}
	if err := json.Unmarshal(b, &element); err != nil {
		return err
	}
	switch item := element.(type) {
	case string:
		*s = String{item}
		return nil
	case map[string]interface{}:
		if file, ok := item["file"]; ok {
			if filename, ok := file.(string); ok {
				data, err := os.ReadFile(filename)
				if err != nil {
					return fmt.Errorf("failed to read referenced file %q: %w", filename, err)
				}
				*s = String{string(data)}
				return nil
			}
		}
	}
	return errInvalidType
}

// Decode implements envconfig.Decoder.
func (s *String) Decode(value string) error {
	*s = String{value}
	return nil
}
