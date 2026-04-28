package xtypes

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"strconv"
	"testing"
)

func TestGroupIDAndUserIDMarshalUnmarshal(t *testing.T) {
	t.Parallel()

	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("current user: %v", err)
	}

	currentUID, err := strconv.Atoi(currentUser.Uid)
	if err != nil {
		t.Fatalf("uid atoi: %v", err)
	}

	currentGID, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		t.Fatalf("gid atoi: %v", err)
	}

	groupID := GroupID(currentGID)
	userID := UserID(currentUID)

	if got := groupID.String(); got == "" {
		t.Fatalf("group string empty")
	}
	if got := userID.String(); got == "" {
		t.Fatalf("user string empty")
	}

	groupJSON, err := (&groupID).MarshalJSON()
	if err != nil {
		t.Fatalf("marshal group: %v", err)
	}
	userJSON, err := (&userID).MarshalJSON()
	if err != nil {
		t.Fatalf("marshal user: %v", err)
	}

	if len(groupJSON) == 0 || len(userJSON) == 0 {
		t.Fatalf("expected non-empty JSON")
	}

	groupText, err := (&groupID).MarshalText()
	if err != nil {
		t.Fatalf("marshal group text: %v", err)
	}
	userText, err := (&userID).MarshalText()
	if err != nil {
		t.Fatalf("marshal user text: %v", err)
	}
	if len(groupText) == 0 || len(userText) == 0 {
		t.Fatalf("expected non-empty text")
	}

	var parsedGroup GroupID
	if err := (&parsedGroup).UnmarshalJSON([]byte(`-1`)); err != nil {
		t.Fatalf("unmarshal group -1: %v", err)
	}
	if parsedGroup != GroupID(os.Getgid()) {
		t.Fatalf("expected os gid, got %d", parsedGroup)
	}

	if err := (&parsedGroup).UnmarshalJSON([]byte(`""`)); err != nil {
		t.Fatalf("unmarshal group empty string: %v", err)
	}

	var parsedUser UserID
	if err := (&parsedUser).UnmarshalJSON([]byte(`-1`)); err != nil {
		t.Fatalf("unmarshal user -1: %v", err)
	}
	if parsedUser != UserID(os.Getuid()) {
		t.Fatalf("expected os uid, got %d", parsedUser)
	}

	if err := (&parsedUser).UnmarshalJSON([]byte(`""`)); err != nil {
		t.Fatalf("unmarshal user empty string: %v", err)
	}

	if err := (&parsedGroup).UnmarshalText([]byte("")); err != nil {
		t.Fatalf("group text empty: %v", err)
	}
	if err := (&parsedUser).UnmarshalText([]byte("")); err != nil {
		t.Fatalf("user text empty: %v", err)
	}
}

func TestAccountIDErrorsAndHelpers(t *testing.T) {
	t.Parallel()

	var g GroupID
	if err := (&g).UnmarshalJSON([]byte(`70000`)); !errors.Is(err, errGroupIDOutOfRange) {
		t.Fatalf("expected group range error, got %v", err)
	}

	var u UserID
	if err := (&u).UnmarshalJSON([]byte(`70000`)); !errors.Is(err, errUserIDOutOfRange) {
		t.Fatalf("expected user range error, got %v", err)
	}

	_, err := parseAccountID("70000", os.Getuid, lookupUserID)
	if !errors.Is(err, errUserGroupIDOutOfRange) {
		t.Fatalf("expected user/group range error, got %v", err)
	}

	if _, err := lookupUserID("this-user-should-not-exist-xtypes"); err == nil {
		t.Fatalf("expected lookupUserID error")
	}
	if _, err := lookupGroupID("this-group-should-not-exist-xtypes"); err == nil {
		t.Fatalf("expected lookupGroupID error")
	}

	current, err := user.Current()
	if err != nil {
		t.Fatalf("current user: %v", err)
	}
	if _, err := parseAccountID(current.Username, os.Getuid, lookupUserID); err != nil {
		t.Fatalf("parse current username: %v", err)
	}
	if _, err := parseAccountID(current.Gid, os.Getgid, lookupGroupID); err != nil {
		t.Fatalf("parse current gid string: %v", err)
	}
	if _, err := parseAccountID("not-a-real-user-for-xtypes", os.Getuid, lookupUserID); err == nil {
		t.Fatalf("expected parseAccountID lookup error")
	}
}

func TestJSONUserIDAndJSONGroupID(t *testing.T) {
	t.Parallel()

	var jsonUser JSONUserID = JSONUserID(UserID(os.Getuid()))
	data, err := json.Marshal(jsonUser)
	if err != nil {
		t.Fatalf("marshal JSONUserID: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected JSONUserID JSON")
	}

	var parsedUser JSONUserID
	if err := json.Unmarshal([]byte(`"-1"`), &parsedUser); err != nil {
		t.Fatalf("unmarshal JSONUserID: %v", err)
	}
	if UserID(parsedUser) != UserID(os.Getuid()) {
		t.Fatalf("expected uid from -1")
	}

	text, err := parsedUser.MarshalText()
	if err != nil {
		t.Fatalf("marshal text JSONUserID: %v", err)
	}
	if err := (&parsedUser).UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal text JSONUserID: %v", err)
	}

	var jsonGroup JSONGroupID = JSONGroupID(GroupID(os.Getgid()))
	data, err = json.Marshal(jsonGroup)
	if err != nil {
		t.Fatalf("marshal JSONGroupID: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected JSONGroupID JSON")
	}

	var parsedGroup JSONGroupID
	if err := json.Unmarshal([]byte(`"-1"`), &parsedGroup); err != nil {
		t.Fatalf("unmarshal JSONGroupID: %v", err)
	}
	if GroupID(parsedGroup) != GroupID(os.Getgid()) {
		t.Fatalf("expected gid from -1")
	}

	text, err = parsedGroup.MarshalText()
	if err != nil {
		t.Fatalf("marshal text JSONGroupID: %v", err)
	}
	if err := (&parsedGroup).UnmarshalText(text); err != nil {
		t.Fatalf("unmarshal text JSONGroupID: %v", err)
	}
}

