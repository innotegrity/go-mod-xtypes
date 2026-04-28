package xtypes

import (
	"context"

	"go.innotegrity.dev/mod/xerrors"
)

const (
	// InvalidParameterErrorCode indicates that an invalid value or type was passed as a parameter to a function.
	InvalidParameterErrorCode = 1

	// PathErrorCode indicates there was a general error while working with the path.
	PathErrorCode = 200

	// PathChmodErrorCode indicates there was an error while changing the permissions of the path.
	PathChmodErrorCode = 201

	// PathChownErrorCode indicates there was an error while changing the ownership of the path.
	PathChownErrorCode = 202

	// PathCreateErrorCode indicates there was an error while creating the path.
	PathCreateErrorCode = 203

	// PathOpenFileErrorCode indicates there was an error while opening the file.
	PathOpenFileErrorCode = 204

	// PathWriteErrorCode indicates there was an error while writing to the file.
	PathWriteErrorCode = 205
)

// InvalidParameterError is returned when a function receives an invalid argument.
type InvalidParameterError struct{ *xerrors.XError }

// PathError is returned for general filesystem path handling failures (stat, resolve, etc.).
type PathError struct{ *xerrors.XError }

// PathChmodError is returned when chmod on a path fails.
type PathChmodError struct{ *xerrors.XError }

// PathChownError is returned when chown on a path fails.
type PathChownError struct{ *xerrors.XError }

// PathCreateError is returned when creating a directory or path fails.
type PathCreateError struct{ *xerrors.XError }

// PathOpenFileError is returned when opening (or preparing to open) a file fails.
type PathOpenFileError struct{ *xerrors.XError }

// PathWriteError is returned when writing to a file fails.
type PathWriteError struct{ *xerrors.XError }

// newPathError creates a new [PathError] with the given error and format.
func newPathError(ctx context.Context, err error, format string, args ...any) *PathError {
	return &PathError{
		XError: xerrors.NewXError(ctx, err, PathErrorCode, format, args...),
	}
	/*
		return xerrors.WrapfAs(func(e *xerrors.XError) *PathError {
			return &PathError{XError: e}
		}, err, PathErrorCode, format, args...)
	*/
}

// newPathChmodError creates a new [PathChmodError] with the given error and format.
func newPathChmodError(ctx context.Context, err error, format string, args ...any) *PathChmodError {
	return &PathChmodError{
		XError: xerrors.NewXError(ctx, err, PathChmodErrorCode, format, args...),
	}
	/*
		return xerrors.WrapfAs(func(e *xerrors.XError) *PathChmodError {
			return &PathChmodError{XError: e}
		}, err, PathChmodErrorCode, format, args...)
	*/
}

// newPathChownError creates a new [PathChownError] with the given error and format.
func newPathChownError(ctx context.Context, err error, format string, args ...any) *PathChownError {
	return &PathChownError{
		XError: xerrors.NewXError(ctx, err, PathChownErrorCode, format, args...),
	}
	/*
		return xerrors.WrapfAs(func(e *xerrors.XError) *PathChownError {
			return &PathChownError{XError: e}
		}, err, PathChownErrorCode, format, args...)
	*/
}

// newPathCreateError creates a new [PathCreateError] with the given error and format.
func newPathCreateError(ctx context.Context, err error, format string, args ...any) *PathCreateError {
	return &PathCreateError{
		XError: xerrors.NewXError(ctx, err, PathCreateErrorCode, format, args...),
	}
	/*
		return xerrors.WrapfAs(func(e *xerrors.XError) *PathCreateError {
			return &PathCreateError{XError: e}
		}, err, PathCreateErrorCode, format, args...)
	*/
}

// newPathOpenFileError creates a new [PathOpenFileError] with the given error and format.
func newPathOpenFileError(ctx context.Context, err error, format string, args ...any) *PathOpenFileError {
	return &PathOpenFileError{
		XError: xerrors.NewXError(ctx, err, PathOpenFileErrorCode, format, args...),
	}
	/*
		return xerrors.WrapfAs(func(e *xerrors.XError) *PathOpenFileError {
			return &PathOpenFileError{XError: e}
		}, err, PathOpenFileErrorCode, format, args...)
	*/
}

// newPathWriteError creates a new [PathWriteError] with the given error and format.
func newPathWriteError(ctx context.Context, err error, format string, args ...any) *PathWriteError {
	return &PathWriteError{
		XError: xerrors.NewXError(ctx, err, PathWriteErrorCode, format, args...),
	}
	/*
		return xerrors.WrapfAs(func(e *xerrors.XError) *PathWriteError {
			return &PathWriteError{XError: e}
		}, err, PathWriteErrorCode, format, args...)
	*/
}
