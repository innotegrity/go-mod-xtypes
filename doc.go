// Package xtypes provides reusable value types and helpers for configuration-driven Go programs: human-friendly
// durations and byte sizes, JSON-friendly file modes and POSIX user/group identifiers, filesystem path operations
// with optional ownership and permissions, generic sets, UUID generation, and small slice utilities. Types
// integrate with encoding/json and encoding/text where noted, so they work well in config structs and APIs.
//
// Import as:
//
//	import "go.innotegrity.dev/mod/xtypes"
//
// # Durations
//
// [Duration] wraps [time.Duration] and adds parsing for calendar-style suffixes beyond the standard library:
// mo (30 days), w (7 days), d (24 hours), and y (365 days), in addition to strings accepted by [time.ParseDuration].
// Empty strings unmarshal to zero. [ParseDuration] parses these forms for use outside JSON/text unmarshaling.
//
// If you need stable custom marshaling when values are carried inside interfaces or map values, use [JSONDuration].
//
// Example:
//
//	d, err := types.ParseDuration("2w")
//	if err != nil {
//		log.Fatal(err)
//	}
//	_ = time.Duration(d) // 336h0m0s
//
// # Sizes
//
// [Size] represents a byte count as float64 with parsing from strings like "512", "1.5MB", or "256KiB".
// Decimal SI prefixes (k, m, g, t, p) use powers of 1000; binary suffixes (KiB, MiB, GiB, TiB, PiB) use powers of 1024.
// [ParseSize] returns the size in bytes. Empty strings parse to zero.
//
// If you need stable custom marshaling when values are carried inside interfaces or map values, use [JSONSize].
//
// Example:
//
//	s, err := types.ParseSize("10MiB")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(s) // human-readable form via String()
//
// # File modes and accounts
//
// [FileMode] stores Unix permission bits and marshals to JSON as octal strings (e.g. "0644"). Use [FileMode.OSFileMode]
// when calling [os.Chmod], [os.MkdirAll], or [os.OpenFile].
//
// [UserID] and [GroupID] represent POSIX uid/gid values. JSON and text unmarshaling accept numeric ids, "-1" for the
// current user or group, empty string for the current user/group, or a system account name resolved via [os/user].
// [UserID.String] and [GroupID.String] resolve numeric ids to names when possible for logging.
//
// If you need stable custom marshaling when values are carried inside interfaces or map values, use [JSONUserID]
// and [JSONGroupID] transport wrappers.
//
// # Paths and errors
//
// [LocalPath] describes a filesystem target with optional directory and file modes, owner, group, and flags to
// auto-create parents, chmod, or chown after create/open. Methods such as [LocalPath.MkdirAll],
// [LocalPath.OpenFile], [LocalPath.WriteFile], [LocalPath.Chmod], [LocalPath.Chown], and [LocalPath.ToAbs] return
// [go.innotegrity.dev/mod/xerrors.Error] values with structured attributes from [LocalPath.Attrs]. Concrete types such
// as [*PathError], [*PathChmodError], and [*PathCreateError] embed
// [*go.innotegrity.dev/mod/xerrors.XError]; stable numeric codes are defined as [PathErrorCode], [PathChmodErrorCode],
// [PathCreateErrorCode], and related constants.
//
// Example:
//
//	p := types.LocalPath{
//		FSPath:           "/var/lib/app/data.json",
//		FileMode:         0o644,
//		DirMode:          0o755,
//		AutoCreateParent: true,
//		AutoChmod:        true,
//	}
//	if xerr := p.WriteFile([]byte("{}"), true); xerr != nil {
//		log.Fatal(xerr)
//	}
//
// # Sets
//
// [Set] is a generic map-backed set for comparable types. Use [NewSet], [Set.Add], [Set.Contains], [Set.Union],
// [Set.Intersection], and [Set.Members] for membership and set algebra.
//
// Example:
//
//	s := types.NewSet("a", "b")
//	s.Add("c")
//	t := types.NewSet("b", "d")
//	u := s.Union(t) // contains a, b, c, d
//	_ = u.Intersection(s)
//
// # UUIDs
//
// [NewUUID] returns an uppercase UUID string and prefers UUID v7 from [github.com/google/uuid]; if v7 fails it
// returns a v8-style identifier or an error if cryptographic random bytes cannot be read for the fallback.
//
// # Utilities
//
// [AnySlice] converts []T to []any for APIs that require []any.
//
// Types that implement [AnyUnmarshaler] can participate in generic unmarshaling pipelines that dispatch on [any].
package xtypes
