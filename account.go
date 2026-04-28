package xtypes

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

var (
	errGroupIDOutOfRange     = errors.New("group ID must be between -1 and 65535, inclusively")
	errUserIDOutOfRange      = errors.New("user ID must be between -1 and 65535, inclusively")
	errUserGroupIDOutOfRange = errors.New("user/group ID must be between -1 and 65535, inclusively")
)

// GroupID represents a Linux or MacOS group ID.
type GroupID int

// UserID represents a Linux or MacOS user ID.
type UserID int

// JSONUserID is a transport type over [UserID] that provides stable JSON marshaling behavior
// when passed as a value (for example in maps or interface fields), without an extra nested field.
//
//nolint:recvcheck // MarshalJSON must be a value receiver for stable transport behavior; UnmarshalJSON must mutate.
type JSONUserID UserID

// JSONGroupID is a transport type over [GroupID] that provides stable JSON marshaling behavior
// when passed as a value (for example in maps or interface fields), without an extra nested field.
//
//nolint:recvcheck // MarshalJSON must be a value receiver for stable transport behavior; UnmarshalJSON must mutate.
type JSONGroupID GroupID

// MarshalJSON marshals the [GroupID] object to JSON.
func (g *GroupID) MarshalJSON() ([]byte, error) {
	result, err := json.Marshal(g.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group ID: %w", err)
	}

	return result, nil
}

// MarshalText marshals the [GroupID] object to plain text.
func (g *GroupID) MarshalText() ([]byte, error) {
	return []byte(g.String()), nil
}

// String returns the [GroupID] object as a string.
func (g *GroupID) String() string {
	gid := "0"
	if g != nil {
		gid = fmt.Sprintf("%d", *g)
	}

	group, err := user.LookupGroupId(gid)
	if err != nil {
		return gid
	}

	return group.Name
}

// UnmarshalJSON parses the JSON data into a [GroupID] object.
//
// If an empty string is supplied, the current group is stored.
func (g *GroupID) UnmarshalJSON(data []byte) error {
	// first see if we have an actual integer value
	var groupID int

	err := json.Unmarshal(data, &groupID)
	if err == nil {
		if groupID < -1 || groupID > 65535 {
			return errGroupIDOutOfRange
		}

		// -1 indicates that we should use the current user/group
		if groupID == -1 {
			groupID = os.Getgid()
		}

		*g = GroupID(groupID)

		return nil
	}

	// try and parse the data as a string
	var strID string

	err = json.Unmarshal(data, &strID)
	if err != nil {
		return fmt.Errorf("failed to unmarshal group ID: %w", err)
	}

	groupID, err = parseAccountID(strID, os.Getgid, lookupGroupID)
	if err != nil {
		return err
	}

	*g = GroupID(groupID)

	return nil
}

// UnmarshalText parses the text into a [GroupID] object.
//
// If an empty string is supplied, the current group is stored.
func (g *GroupID) UnmarshalText(data []byte) error {
	id, err := parseAccountID(string(data), os.Getgid, lookupGroupID)
	if err != nil {
		return err
	}

	*g = GroupID(id)

	return nil
}

// MarshalJSON marshals the [UserID] object to JSON.
func (u *UserID) MarshalJSON() ([]byte, error) {
	result, err := json.Marshal(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user ID: %w", err)
	}

	return result, nil
}

// MarshalText marshals the [UserID] object to plain text.
func (u *UserID) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// String returns the [UserID] object as a string.
func (u *UserID) String() string {
	uid := "0"
	if u != nil {
		uid = fmt.Sprintf("%d", *u)
	}

	user, err := user.LookupId(uid)
	if err != nil {
		return uid
	}

	return user.Username
}

// UnmarshalJSON parses the JSON data into a [UserID] object.
//
// If an empty string is supplied, the current user is stored.
func (u *UserID) UnmarshalJSON(data []byte) error {
	// first see if we have an actual integer value
	var userID int

	err := json.Unmarshal(data, &userID)
	if err == nil {
		if userID < -1 || userID > 65535 {
			return errUserIDOutOfRange
		}

		// -1 indicates that we should use the current user/group
		if userID == -1 {
			userID = os.Getuid()
		}

		*u = UserID(userID)

		return nil
	}

	// try and parse the data as a string
	var strID string

	err = json.Unmarshal(data, &strID)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user ID: %w", err)
	}

	userID, err = parseAccountID(strID, os.Getuid, lookupUserID)
	if err != nil {
		return err
	}

	*u = UserID(userID)

	return nil
}

// UnmarshalText parses the text into a [UserID] object.
//
// If an empty string is supplied, the current user is stored.
func (u *UserID) UnmarshalText(data []byte) error {
	id, err := parseAccountID(string(data), os.Getuid, lookupUserID)
	if err != nil {
		return err
	}

	*u = UserID(id)

	return nil
}

// MarshalJSON marshals the [JSONUserID] using [UserID] JSON behavior.
func (j JSONUserID) MarshalJSON() ([]byte, error) {
	userID := UserID(j)

	data, err := (&userID).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONUserID: %w", err)
	}

	return data, nil
}

// MarshalText marshals [JSONUserID] as plain text.
func (j JSONUserID) MarshalText() ([]byte, error) {
	userID := UserID(j)

	data, err := (&userID).MarshalText()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal text JSONUserID: %w", err)
	}

	return data, nil
}

// UnmarshalJSON unmarshals JSON into [JSONUserID].
func (j *JSONUserID) UnmarshalJSON(data []byte) error {
	var userID UserID

	err := (&userID).UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONUserID: %w", err)
	}

	*j = JSONUserID(userID)

	return nil
}

// UnmarshalText unmarshals text into [JSONUserID].
func (j *JSONUserID) UnmarshalText(data []byte) error {
	var userID UserID

	err := (&userID).UnmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text JSONUserID: %w", err)
	}

	*j = JSONUserID(userID)

	return nil
}

// MarshalJSON marshals the [JSONGroupID] using [GroupID] JSON behavior.
func (j JSONGroupID) MarshalJSON() ([]byte, error) {
	groupID := GroupID(j)

	data, err := (&groupID).MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSONGroupID: %w", err)
	}

	return data, nil
}

// MarshalText marshals [JSONGroupID] as plain text.
func (j JSONGroupID) MarshalText() ([]byte, error) {
	groupID := GroupID(j)

	data, err := (&groupID).MarshalText()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal text JSONGroupID: %w", err)
	}

	return data, nil
}

// UnmarshalJSON unmarshals JSON into [JSONGroupID].
func (j *JSONGroupID) UnmarshalJSON(data []byte) error {
	var groupID GroupID

	err := (&groupID).UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONGroupID: %w", err)
	}

	*j = JSONGroupID(groupID)

	return nil
}

// UnmarshalText unmarshals text into [JSONGroupID].
func (j *JSONGroupID) UnmarshalText(data []byte) error {
	var groupID GroupID

	err := (&groupID).UnmarshalText(data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal text JSONGroupID: %w", err)
	}

	*j = JSONGroupID(groupID)

	return nil
}

// lookupGroupID attempts to lookup the ID of the given group.
func lookupGroupID(name string) (string, error) {
	g, err := user.LookupGroup(name)
	if err != nil {
		return "", fmt.Errorf("failed to lookup group named '%s': %w", name, err)
	}

	return g.Gid, nil
}

// lookupUserID attempts to lookup the ID of the given user.
func lookupUserID(name string) (string, error) {
	u, err := user.Lookup(name)
	if err != nil {
		return "", fmt.Errorf("failed to lookup user named '%s': %w", name, err)
	}

	return u.Uid, nil
}

// parseAccountID handles parsing the given data into a user or group ID.
func parseAccountID(data string, getCurrentID func() int, lookupAccount func(string) (string, error)) (int, error) {
	// empty string indicates that we should use the current user/group
	var accountID int
	if data == "" {
		accountID = getCurrentID()

		return accountID, nil
	}

	// try and convert the string to an integer
	accountID, err := strconv.Atoi(data)
	if err == nil {
		if accountID < -1 || accountID > 65535 {
			return -2, errUserGroupIDOutOfRange
		}

		// -1 indicates that we should use the current user/group
		if accountID == -1 {
			accountID = getCurrentID()
		}

		return accountID, nil
	}

	// try and look up the user/group
	strID, err := lookupAccount(data)
	if err != nil {
		return -2, err
	}

	accountID, err = strconv.Atoi(strID)
	if err != nil {
		return -2, fmt.Errorf("failed to convert user/group ID '%s' to an integer: %w", data, err)
	}

	return accountID, nil
}
