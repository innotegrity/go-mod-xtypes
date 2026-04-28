package xtypes

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	fileModeParseBase = 0
	fileModeParseBits = 32
	fileModeFormat    = "%#o"
)

// FileMode represents a file or directory mode.
type FileMode int

// JSONFileMode is a transport type over [FileMode] that provides stable JSON marshaling behavior
// when passed as a value (for example in maps or interface fields), without an extra nested field.
//
//nolint:recvcheck // MarshalJSON must be a value receiver for stable transport behavior; UnmarshalJSON must mutate.
type JSONFileMode FileMode

// MarshalJSON marshals the [FileMode] object to JSON.
func (m *FileMode) MarshalJSON() ([]byte, error) {
	result, err := json.Marshal(m.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal file mode: %w", err)
	}

	return result, nil
}

// MarshalText marshals the [FileMode] object to plain text.
func (m *FileMode) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// OSFileMode returns the [os.FileMode] equivalent of the object.
func (m *FileMode) OSFileMode() os.FileMode {
	if m == nil || *m < 0 {
		return os.FileMode(0)
	}

	modeValue := int64(*m)
	if modeValue > math.MaxUint32 {
		return os.FileMode(0)
	}

	return os.FileMode(uint32(modeValue))
}

// String returns the [FileMode] object as a string.
func (m *FileMode) String() string {
	modeValue := FileMode(0)
	if m != nil {
		modeValue = *m
	}

	return fmt.Sprintf(fileModeFormat, modeValue)
}

// UnmarshalJSON parses the JSON data into a [FileMode] object.
func (m *FileMode) UnmarshalJSON(data []byte) error {
	var mode int32

	err := json.Unmarshal(data, &mode)
	if err == nil {
		*m = FileMode(mode)

		return nil
	}

	var modeText string

	err = json.Unmarshal(data, &modeText)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file mode: %w", err)
	}

	return m.UnmarshalText([]byte(modeText))
}

// UnmarshalText parses the text into a [FileMode] object.
func (m *FileMode) UnmarshalText(data []byte) error {
	mode, err := strconv.ParseInt(string(data), fileModeParseBase, fileModeParseBits)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file mode text: %w", err)
	}

	*m = FileMode(mode)

	return nil
}

// MarshalJSON marshals [JSONFileMode] using [FileMode] JSON behavior.
func (m JSONFileMode) MarshalJSON() ([]byte, error) {
	fileModeValue := FileMode(m)

	result, err := (&fileModeValue).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONFileMode: %w", err)
	}

	return result, nil
}

// MarshalText marshals [JSONFileMode] as plain text.
func (m JSONFileMode) MarshalText() ([]byte, error) {
	fileModeValue := FileMode(m)

	result, err := (&fileModeValue).MarshalText()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal text JSONFileMode: %w", err)
	}

	return result, nil
}

// UnmarshalJSON unmarshals JSON into [JSONFileMode].
func (m *JSONFileMode) UnmarshalJSON(data []byte) error {
	var fileModeValue FileMode

	err := (&fileModeValue).UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONFileMode: %w", err)
	}

	*m = JSONFileMode(fileModeValue)

	return nil
}

// UnmarshalText unmarshals text into [JSONFileMode].
func (m *JSONFileMode) UnmarshalText(data []byte) error {
	var fileModeValue FileMode

	err := (&fileModeValue).UnmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text JSONFileMode: %w", err)
	}

	*m = JSONFileMode(fileModeValue)

	return nil
}
