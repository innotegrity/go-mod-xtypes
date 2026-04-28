package xtypes

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"go.innotegrity.dev/mod/xerrors"
)

// LocalPath holds settings for a particular file or folder.
type LocalPath struct {
	// AutoChmod indicates if the permissions of the file or directory should be changed when creating or opening it.
	AutoChmod bool `json:"autoChmod" mapstructure:"autoChmod" yaml:"autoChmod"`

	// AutoChown indicates if the ownership of the file or directory should be changed when creating or opening it.
	AutoChown bool `json:"autoChown" mapstructure:"autoChown" yaml:"autoChown"`

	// AutoCreateParent indicates if any parent folders should be created if they do not exist when creating oropening
	// a file.
	AutoCreateParent bool `json:"autoCreateParent" mapstructure:"autoCreateParent" yaml:"autoCreateParent"`

	// DirMode is the mode that should be used when creating the directory or any parent directories.
	DirMode FileMode `json:"dirMode" mapstructure:"dirMode" yaml:"dirMode"`

	// FileMode is the mode that should be used when creating the file.
	FileMode FileMode `json:"fileMode" mapstructure:"fileMode" yaml:"fileMode"`

	// FSPath is the path to the file or directory on the filesystem.
	FSPath string `json:"path" mapstructure:"path" yaml:"path"`

	// Group is the group name or ID that should own the file or directory.
	Group GroupID `json:"group" mapstructure:"group" yaml:"group"`

	// Owner is the user name or ID that should own the file or directory.
	Owner UserID `json:"owner" mapstructure:"owner" yaml:"owner"`
}

// ToAbs attempts to convert the filesystem path to an absolute path using [context.Background].
func (p *LocalPath) ToAbs() xerrors.Error {
	return p.ToAbsContext(context.Background())
}

// ToAbsContext attempts to convert the filesystem path to an absolute path.
//
// This function modifies the [LocalPath.FSPath] field in place.
//
// The given context is used to transport [xerrors.Error] options (for example caller capture)
// into any returned structured errors.
//
// This function may return any of the following errors:
//   - [*PathError]: there was a general error while working with the path
func (p *LocalPath) ToAbsContext(ctx context.Context) xerrors.Error {
	path, err := filepath.Abs(p.FSPath)
	if err != nil {
		return newPathError(ctx, err, "failed to convert '%s' to an absolute path: %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path": p.FSPath,
			})
	}

	p.FSPath = path

	return nil
}

// Attrs returns the attributes of the path which can be attached to errors or log messages.
func (p *LocalPath) Attrs() map[string]any {
	return map[string]any{
		"dir_mode":  fmt.Sprintf("%o", p.DirMode),
		"file_mode": fmt.Sprintf("%o", p.FileMode),
		"group":     p.Group.String(),
		"owner":     p.Owner.String(),
		"path":      p.FSPath,
	}
}

// Chmod sets the stored permissions on the path using [context.Background].
func (p *LocalPath) Chmod() xerrors.Error {
	return p.ChmodContext(context.Background())
}

// ChmodContext sets the stored permissions on the path.
//
// The given context is used to transport [xerrors.Error] options (for example caller capture)
// into any returned structured errors.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the file/folder
//   - [*PathError]: there was a general error while working with the path
func (p *LocalPath) ChmodContext(ctx context.Context) xerrors.Error {
	pathInfo, err := os.Stat(p.FSPath)
	if err != nil {
		return newPathError(ctx, err, "failed to change permissions of '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path": p.FSPath,
			})
	}

	mode := p.FileMode
	if pathInfo.IsDir() {
		mode = p.DirMode
	}

	err = os.Chmod(p.FSPath, mode.OSFileMode())
	if err != nil {
		return newPathChmodError(ctx, err, "failed to change permissions of '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path":     p.FSPath,
				"new_mode": fmt.Sprintf("%#o", mode),
			})
	}

	return nil
}

// Chown sets the stored ownership for the path using [context.Background].
func (p *LocalPath) Chown() xerrors.Error {
	return p.ChownContext(context.Background())
}

// ChownContext sets the stored ownership for the path.
//
// The given context is used to transport [xerrors.Error] options (for example caller capture)
// into any returned structured errors.
//
// This function may return any of the following errors:
//   - [*PathChownError]: there was an error while changing ownership of the file/folder
func (p *LocalPath) ChownContext(ctx context.Context) xerrors.Error {
	// only works for root
	if os.Geteuid() != 0 {
		return nil
	}

	err := os.Chown(p.FSPath, int(p.Owner), int(p.Group))
	if err != nil {
		return newPathChownError(ctx, err, "failed to change ownership of '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path":      p.FSPath,
				"new_owner": p.Owner.String(),
				"new_group": p.Group.String(),
			})
	}

	return nil
}

// MkdirAll creates the given path and any parent folders if they do not exist using [context.Background].
func (p *LocalPath) MkdirAll() xerrors.Error {
	return p.MkdirAllContext(context.Background())
}

// MkdirAllContext creates the given path and any parent folders if they do not exist.
//
// If [LocalPath.AutoChmod] is true, the permissions will be set to the [LocalPath.DirMode] value.
// If [LocalPath.AutoChown] is true, the ownership will be set to the [LocalPath.Owner] and [LocalPath.Group] values.
//
// The given context is used to transport [xerrors.Error] options (for example caller capture)
// into any returned structured errors.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the folder
//   - [*PathChownError]: there was an error while changing ownership of the folder
//   - [*PathCreateError]: there was an error while creating the folder
//   - [*PathError]: there was a general error while working with the path
func (p *LocalPath) MkdirAllContext(ctx context.Context) xerrors.Error {
	// create the folder
	err := os.MkdirAll(p.FSPath, p.DirMode.OSFileMode())
	if err != nil {
		return newPathCreateError(ctx, err, "failed to create path '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path":     p.FSPath,
				"dir_mode": fmt.Sprintf("%o", p.DirMode),
			})
	}

	// set ownership and permissions
	if p.AutoChmod {
		xerr := p.ChmodContext(ctx)
		if xerr != nil {
			return xerr
		}
	}

	if p.AutoChown {
		xerr := p.ChownContext(ctx)
		if xerr != nil {
			return xerr
		}
	}

	return nil
}

// OpenFile creates/opens the file and returns its handle using [context.Background].
func (p *LocalPath) OpenFile(flags int) (*os.File, xerrors.Error) {
	return p.OpenFileContext(context.Background(), flags)
}

// OpenFileContext creates/opens the file and returns its handle.
//
// If [LocalPath.AutoCreateParent] is true, [LocalPath.MkdirAllContext] will be called on the file's parent folder
// first. If [LocalPath.AutoChmod] is true, the permissions will be set to the [LocalPath.DirMode] value.
// If [LocalPath.AutoChown] is true, the ownership will be set to the [LocalPath.Owner] and [LocalPath.Group] values.
//
// The given context is used to transport [xerrors.Error] options (for example caller capture)
// into any returned structured errors.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the file/parent folder
//   - [*PathChownError]: there was an error while changing ownership of the file/parent folder
//   - [*PathCreateError]: there was an error while creating the parent folder
//   - [*PathError]: there was a general error while working with the path
//   - [*PathOpenFileError]: there was an error while opening the file
func (p *LocalPath) OpenFileContext(ctx context.Context, flags int) (*os.File, xerrors.Error) {
	// create parent folder if desired
	if p.AutoCreateParent {
		parent := LocalPath{
			DirMode: p.DirMode,
			Group:   p.Group,
			Owner:   p.Owner,
			FSPath:  path.Dir(p.FSPath),
		}

		xerr := parent.MkdirAllContext(ctx)
		if xerr != nil {
			return nil, newPathOpenFileError(ctx, xerr, "failed to open file '%s': %s", p.FSPath,
				xerr.Error()).WithAttrs(map[string]any{
				"file":      p.FSPath,
				"file_mode": fmt.Sprintf("%o", p.FileMode),
			})
		}
	}

	// open the file
	file, err := os.OpenFile(p.FSPath, flags, p.FileMode.OSFileMode())
	if err != nil {
		return nil, newPathOpenFileError(ctx, err, "failed to open file '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"file":      p.FSPath,
				"file_mode": fmt.Sprintf("%o", p.FileMode),
			})
	}

	// set ownership and permissions
	if p.AutoChmod {
		xerr := p.ChmodContext(ctx)
		if xerr != nil {
			_ = file.Close()

			return nil, xerr
		}
	}

	if p.AutoChown {
		xerr := p.ChownContext(ctx)
		if xerr != nil {
			_ = file.Close()

			return nil, xerr
		}
	}

	return file, nil
}

// WriteFile writes the given data to the file using [context.Background].
func (p *LocalPath) WriteFile(data []byte, overwrite bool) xerrors.Error {
	return p.WriteFileContext(context.Background(), data, overwrite)
}

// WriteFileContext writes the given data to the file.
//
// This function uses the [LocalPath.OpenFileContext] function to create/open the file before writing to it. It
// automatically closes the file after writing to it.
//
// The given context is used to transport [xerrors.Error] options (for example caller capture)
// into any returned structured errors.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the file/parent folder
//   - [*PathChownError]: there was an error while changing ownership of the file/parent folder
//   - [*PathCreateError]: there was an error while creating the parent folder
//   - [*PathError]: there was a general error while working with the path
//   - [*PathOpenFileError]: there was an error while opening the file
//   - [*PathWriteError]: there was an error while writing to the file
func (p *LocalPath) WriteFileContext(ctx context.Context, data []byte, overwrite bool) xerrors.Error {
	flags := os.O_CREATE | os.O_RDWR
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}

	handle, xerr := p.OpenFileContext(ctx, flags)
	if xerr != nil {
		return xerr
	}

	defer func() {
		_ = handle.Close()
	}()

	_, err := handle.Write(data)
	if err != nil {
		return newPathWriteError(ctx, err, "failed to write to file '%s': %s", p.FSPath, err.Error()).
			WithAttr("file", p.FSPath)
	}

	return nil
}
