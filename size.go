package xtypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Decimal (SI) byte multipliers — powers of 1000.
const (
	sizeKilo = 1000
	sizeMega = sizeKilo * sizeKilo
	sizeGiga = sizeMega * sizeKilo
	sizeTera = sizeGiga * sizeKilo
	sizePeta = sizeTera * sizeKilo
)

// Binary (IEC) byte multipliers — powers of 1024.
const (
	sizeKibi = 1024
	sizeMebi = sizeKibi * sizeKibi
	sizeGibi = sizeMebi * sizeKibi
	sizeTebi = sizeGibi * sizeKibi
	sizePebi = sizeTebi * sizeKibi
)

const (
	sizeParseBase    = 10
	sizeParseFloat64 = 64
	sizeParseInt64   = 64
	sizeFormatBytes  = "%g bytes"
	sizeFormatKB     = "%gKB"
	sizeFormatMB     = "%gMB"
	sizeFormatGB     = "%gGB"
	sizeFormatTB     = "%gTB"
	sizeFormatPB     = "%gPB"
)

var errSizeOverflow = errors.New("size exceeds maximum size of a 64-bit integer")

// Size is an extended version of a float64 which allows abbreviating sizes by adding a suffix.
//
// A value without a suffix must be an integer and will be treated as a size in bytes. You may also add one of
// the following suffixes after the value to indicate a different measurement:
//
//	b | bytes = size is in bytes
//	k | kb = size is in kilobytes (where 1k = 1000^1 bytes)
//	kib = size is in kibibytes (where 1kib = 1024^1 bytes)
//	m | mb = size is in megabytes (where 1m = 1000^2 bytes)
//	mib = size is in mebibytes (where 1mib = 1024^2 bytes)
//	g | gb = size is in gigabytes (where 1g = 1000^3 bytes)
//	gib = size is in gibibytes (where 1gib = 1024^3 bytes)
//	t | tb = size is in terabytes (where 1g = 1000^4 bytes)
//	tib = size is in tebibytes (where 1gib = 1024^4 bytes)
//	p | pb = size is in petabytes (where 1p = 1000^5 bytes)
//	pib = size is in pebibytes (where 1pib = 1024^5 bytes)
type Size float64

// JSONSize is a transport type over [Size] that provides stable JSON marshaling behavior
// when passed as a value (for example in maps or interface fields), without an extra nested field.
//
//nolint:recvcheck // MarshalJSON must be a value receiver for stable transport behavior; UnmarshalJSON must mutate.
type JSONSize Size

// scaleSize multiplies a parsed scalar by a byte-unit multiplier with overflow checking.
func scaleSize(f, bytesPerUnit float64) (Size, error) {
	if f > math.MaxFloat64/bytesPerUnit {
		return 0, errSizeOverflow
	}

	return Size(f * bytesPerUnit), nil
}

// ParseSize parses the given string into a [Size] object.
//
// If an empty string is supplied, 0 is returned.
func ParseSize(size string) (Size, error) {
	if size == "" {
		return 0, nil
	}

	sizePattern := regexp.MustCompile(
		`^(\d*\.\d+|\d+\.\d*|\d+)(\s*?)(?i)(b|bytes|k|kb|kib|m|mb|mib|g|gb|gib|t|tb|tib|p|pb|pib)$`)

	matches := sizePattern.FindStringSubmatch(size)
	if matches == nil {
		ival64, err := strconv.ParseInt(size, sizeParseBase, sizeParseInt64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse size '%s': %w", size, err)
		}

		return Size(ival64), nil
	}

	fval64, err := strconv.ParseFloat(matches[1], sizeParseFloat64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse size '%s': %w", size, err)
	}

	switch strings.ToLower(matches[3]) {
	case "k", "kb":
		return scaleSize(fval64, sizeKilo)
	case "kib":
		return scaleSize(fval64, sizeKibi)
	case "m", "mb":
		return scaleSize(fval64, sizeMega)
	case "mib":
		return scaleSize(fval64, sizeMebi)
	case "g", "gb":
		return scaleSize(fval64, sizeGiga)
	case "gib":
		return scaleSize(fval64, sizeGibi)
	case "t", "tb":
		return scaleSize(fval64, sizeTera)
	case "tib":
		return scaleSize(fval64, sizeTebi)
	case "p", "pb":
		return scaleSize(fval64, sizePeta)
	case "pib":
		return scaleSize(fval64, sizePebi)
	default:
		return Size(fval64), nil
	}
}

// MarshalJSON marshals the [Size] object to JSON.
func (s *Size) MarshalJSON() ([]byte, error) {
	result, err := json.Marshal(s.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal size: %w", err)
	}

	return result, nil
}

// MarshalText marshals the [Size] object to plain text.
func (s *Size) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// String returns the [Size] object as a string.
func (s *Size) String() string {
	byteSize := Size(0)
	if s != nil {
		byteSize = *s
	}

	if byteSize < sizeKilo {
		return fmt.Sprintf(sizeFormatBytes, byteSize)
	}

	if byteSize < sizeMega {
		return fmt.Sprintf(sizeFormatKB, float64(byteSize)/float64(sizeKilo))
	}

	if byteSize < sizeGiga {
		return fmt.Sprintf(sizeFormatMB, float64(byteSize)/float64(sizeMega))
	}

	if byteSize < sizeTera {
		return fmt.Sprintf(sizeFormatGB, float64(byteSize)/float64(sizeGiga))
	}

	if byteSize < sizePeta {
		return fmt.Sprintf(sizeFormatTB, float64(byteSize)/float64(sizeTera))
	}

	return fmt.Sprintf(sizeFormatPB, float64(byteSize)/float64(sizePeta))
}

// UnmarshalJSON parses the JSON data into a [Size] object.
//
// If an empty string is supplied, 0 is stored.
func (s *Size) UnmarshalJSON(data []byte) error {
	var fval64 float64

	err := json.Unmarshal(data, &fval64)
	if err == nil {
		*s = Size(fval64)

		return nil
	}

	var sval string

	err = json.Unmarshal(data, &sval)
	if err != nil {
		return fmt.Errorf("failed to unmarshal size: %w", err)
	}

	size, err := ParseSize(sval)
	if err != nil {
		return err
	}

	*s = size

	return nil
}

// UnmarshalText parses the text into a [Size] object.
//
// If an empty string is supplied, 0 is stored.
func (s *Size) UnmarshalText(data []byte) error {
	size, err := ParseSize(string(data))
	if err != nil {
		return err
	}

	*s = size

	return nil
}

// MarshalJSON marshals [JSONSize] using [Size] JSON behavior.
func (s JSONSize) MarshalJSON() ([]byte, error) {
	sizeValue := Size(s)

	result, err := (&sizeValue).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONSize: %w", err)
	}

	return result, nil
}

// MarshalText marshals [JSONSize] as plain text.
func (s JSONSize) MarshalText() ([]byte, error) {
	sizeValue := Size(s)

	result, err := (&sizeValue).MarshalText()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal text JSONSize: %w", err)
	}

	return result, nil
}

// UnmarshalJSON unmarshals JSON into [JSONSize].
func (s *JSONSize) UnmarshalJSON(data []byte) error {
	var sizeValue Size

	err := (&sizeValue).UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONSize: %w", err)
	}

	*s = JSONSize(sizeValue)

	return nil
}

// UnmarshalText unmarshals text into [JSONSize].
func (s *JSONSize) UnmarshalText(data []byte) error {
	var sizeValue Size

	err := (&sizeValue).UnmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text JSONSize: %w", err)
	}

	*s = JSONSize(sizeValue)

	return nil
}
