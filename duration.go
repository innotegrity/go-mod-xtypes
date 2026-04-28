package xtypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	durationParseBase   = 10
	durationParseBits   = 64
	durationHoursInDay  = 24
	durationDaysInWeek  = 7
	durationDaysInMonth = 30
	durationDaysInYear  = 365
	durationSuffixMonth = "mo"
	durationSuffixWeek  = "w"
	durationSuffixDay   = "d"
	durationSuffixYear  = "y"
	durationFormatHours = "%dh"
)

var errDurationOverflow = errors.New("duration exceeds maximum size of a 64-bit integer")

// Duration is an extended version of [time.Duration] which adds the ability to include new durations such as
// months, weeks, days and years using the suffixes "mo", "w", "d" and "y" respectively.
//
// A month is defined as 30 days, a week is defined as 7 days, a day is defined as 24 hours and a year is
// defined as 365 days.
//
// It also supports unmarshaling empty strings to an empty [Duration] object.
type Duration time.Duration

// JSONDuration is a transport type over [Duration] that provides stable JSON marshaling behavior
// when passed as a value (for example in maps or interface fields), without an extra nested field.
//
//nolint:recvcheck // MarshalJSON must be a value receiver for stable transport behavior; UnmarshalJSON must mutate.
type JSONDuration Duration

// ParseDuration parses a string into a [Duration], supporting "mo", "w", "d", and "y" suffixes.
func ParseDuration(dur string) (Duration, error) {
	if dur == "" {
		return Duration(0), nil
	}

	var hourPeriod int64

	switch {
	case strings.HasSuffix(dur, durationSuffixMonth):
		dur = strings.TrimSuffix(dur, durationSuffixMonth)
		hourPeriod = int64(durationDaysInMonth * durationHoursInDay)
	case strings.HasSuffix(dur, durationSuffixWeek):
		dur = strings.TrimSuffix(dur, durationSuffixWeek)
		hourPeriod = int64(durationDaysInWeek * durationHoursInDay)
	case strings.HasSuffix(dur, durationSuffixDay):
		dur = strings.TrimSuffix(dur, durationSuffixDay)
		hourPeriod = int64(durationHoursInDay)
	case strings.HasSuffix(dur, durationSuffixYear):
		dur = strings.TrimSuffix(dur, durationSuffixYear)
		hourPeriod = int64(durationDaysInYear * durationHoursInDay)
	}

	if hourPeriod > 0 {
		val, err := strconv.ParseInt(dur, durationParseBase, durationParseBits)
		if err != nil {
			return 0, fmt.Errorf("failed to parse integer portion of duration '%s': %w", dur, err)
		}

		if val > math.MaxInt64/hourPeriod {
			return 0, fmt.Errorf("duration '%s': %w", dur, errDurationOverflow)
		}

		dur = fmt.Sprintf(durationFormatHours, val*hourPeriod)
	}

	parsedDuration, err := time.ParseDuration(dur)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration '%s': %w", dur, err)
	}

	return Duration(parsedDuration), nil
}

// MarshalJSON marshals the [Duration] object to JSON.
func (d *Duration) MarshalJSON() ([]byte, error) {
	result, err := json.Marshal(d.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal duration: %w", err)
	}

	return result, nil
}

// MarshalText marshals the [Duration] object to plain text.
func (d *Duration) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// String returns the [Duration] object as a string.
func (d *Duration) String() string {
	durationValue := Duration(0)
	if d != nil {
		durationValue = *d
	}

	return time.Duration(durationValue).String()
}

// UnmarshalJSON parses the JSON data into a [Duration] object.
//
// If an empty string is supplied, 0 is stored.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var ival64 int64

	err := json.Unmarshal(data, &ival64)
	if err == nil {
		*d = Duration(ival64)

		return nil
	}

	var sval string

	err = json.Unmarshal(data, &sval)
	if err != nil {
		return fmt.Errorf("failed to unmarshal duration: %w", err)
	}

	dur, err := ParseDuration(sval)
	if err != nil {
		return err
	}

	*d = dur

	return nil
}

// UnmarshalText parses the text into a [Duration] object.
//
// If an empty string is supplied, 0 is stored.
func (d *Duration) UnmarshalText(data []byte) error {
	dur, err := ParseDuration(string(data))
	if err != nil {
		return err
	}

	*d = dur

	return nil
}

// MarshalJSON marshals [JSONDuration] using [Duration] JSON behavior.
func (d JSONDuration) MarshalJSON() ([]byte, error) {
	durationValue := Duration(d)

	result, err := (&durationValue).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONDuration: %w", err)
	}

	return result, nil
}

// MarshalText marshals [JSONDuration] as plain text.
func (d JSONDuration) MarshalText() ([]byte, error) {
	durationValue := Duration(d)

	result, err := (&durationValue).MarshalText()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal text JSONDuration: %w", err)
	}

	return result, nil
}

// UnmarshalJSON unmarshals JSON into [JSONDuration].
func (d *JSONDuration) UnmarshalJSON(data []byte) error {
	var durationValue Duration

	err := (&durationValue).UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONDuration: %w", err)
	}

	*d = JSONDuration(durationValue)

	return nil
}

// UnmarshalText unmarshals text into [JSONDuration].
func (d *JSONDuration) UnmarshalText(data []byte) error {
	var durationValue Duration

	err := (&durationValue).UnmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text JSONDuration: %w", err)
	}

	*d = JSONDuration(durationValue)

	return nil
}
