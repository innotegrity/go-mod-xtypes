package xtypes

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"go.innotegrity.dev/mod/xerrors"
)

var (
	pathAbsFn     = filepath.Abs
	pathStatFn    = os.Stat
	pathChmodFn   = os.Chmod
	pathGeteuidFn = os.Geteuid
	pathChownFn   = os.Chown
	pathMkdirAll  = os.MkdirAll
	pathOpenFile  = os.OpenFile
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

// ToAbs attempts to convert the filesystem path to an absolute path.
//
// This function modifies the [LocalPath.FSPath] field in place.
//
// This function may return any of the following errors:
//   - [*PathError]: there was a general error while working with the path
func (p *LocalPath) ToAbs() xerrors.Error {
	path, err := pathAbsFn(p.FSPath)
	if err != nil {
		return newPathError(err, "failed to convert '%s' to an absolute path: %s", p.FSPath, err.Error()).
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

// Chmod sets the stored permissions on the path.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the file/folder
//   - [*PathError]: there was a general error while working with the path
func (p *LocalPath) Chmod() xerrors.Error {
	pathInfo, err := pathStatFn(p.FSPath)
	if err != nil {
		return newPathError(err, "failed to change permissions of '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path": p.FSPath,
			})
	}

	mode := p.FileMode
	if pathInfo.IsDir() {
		mode = p.DirMode
	}

	err = pathChmodFn(p.FSPath, mode.OSFileMode())
	if err != nil {
		return newPathChmodError(err, "failed to change permissions of '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path":     p.FSPath,
				"new_mode": fmt.Sprintf("%#o", mode),
			})
	}

	return nil
}

// Chown sets the stored ownership for the path.
//
// This function may return any of the following errors:
//   - [*PathChownError]: there was an error while changing ownership of the file/folder
func (p *LocalPath) Chown() xerrors.Error {
	// only works for root
	if pathGeteuidFn() != 0 {
		return nil
	}

	err := pathChownFn(p.FSPath, int(p.Owner), int(p.Group))
	if err != nil {
		return newPathChownError(err, "failed to change ownership of '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path":      p.FSPath,
				"new_owner": p.Owner.String(),
				"new_group": p.Group.String(),
			})
	}

	return nil
}

// MkdirAll creates the given path and any parent folders if they do not exist.
//
// If [LocalPath.AutoChmod] is true, the permissions will be set to the [LocalPath.DirMode] value.
// If [LocalPath.AutoChown] is true, the ownership will be set to the [LocalPath.Owner] and [LocalPath.Group] values.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the folder
//   - [*PathChownError]: there was an error while changing ownership of the folder
//   - [*PathCreateError]: there was an error while creating the folder
//   - [*PathError]: there was a general error while working with the path
func (p *LocalPath) MkdirAll() xerrors.Error {
	// create the folder
	err := pathMkdirAll(p.FSPath, p.DirMode.OSFileMode())
	if err != nil {
		return newPathCreateError(err, "failed to create path '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"path":     p.FSPath,
				"dir_mode": fmt.Sprintf("%o", p.DirMode),
			})
	}

	// set ownership and permissions
	if p.AutoChmod {
		xerr := p.Chmod()
		if xerr != nil {
			return xerr
		}
	}

	if p.AutoChown {
		xerr := p.Chown()
		if xerr != nil {
			return xerr
		}
	}

	return nil
}

// OpenFile creates/opens the file and returns its handle.
//
// If [LocalPath.AutoCreateParent] is true, [LocalPath.MkdirAll] will be called on the file's parent folder first.
// If [LocalPath.AutoChmod] is true, the permissions will be set to the [LocalPath.DirMode] value.
// If [LocalPath.AutoChown] is true, the ownership will be set to the [LocalPath.Owner] and [LocalPath.Group] values.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the file/parent folder
//   - [*PathChownError]: there was an error while changing ownership of the file/parent folder
//   - [*PathCreateError]: there was an error while creating the parent folder
//   - [*PathError]: there was a general error while working with the path
//   - [*PathOpenFileError]: there was an error while opening the file
func (p *LocalPath) OpenFile(flags int) (*os.File, xerrors.Error) {
	// create parent folder if desired
	if p.AutoCreateParent {
		parent := LocalPath{
			DirMode: p.DirMode,
			Group:   p.Group,
			Owner:   p.Owner,
			FSPath:  path.Dir(p.FSPath),
		}

		xerr := parent.MkdirAll()
		if xerr != nil {
			return nil, newPathOpenFileError(xerr, "failed to open file '%s': %s", p.FSPath,
				xerr.Error()).WithAttrs(map[string]any{
				"file":      p.FSPath,
				"file_mode": fmt.Sprintf("%o", p.FileMode),
			})
		}
	}

	// open the file
	file, err := pathOpenFile(p.FSPath, flags, p.FileMode.OSFileMode())
	if err != nil {
		return nil, newPathOpenFileError(err, "failed to open file '%s': %s", p.FSPath, err.Error()).
			WithAttrs(map[string]any{
				"file":      p.FSPath,
				"file_mode": fmt.Sprintf("%o", p.FileMode),
			})
	}

	// set ownership and permissions
	if p.AutoChmod {
		xerr := p.Chmod()
		if xerr != nil {
			_ = file.Close()

			return nil, xerr
		}
	}

	if p.AutoChown {
		xerr := p.Chown()
		if xerr != nil {
			_ = file.Close()

			return nil, xerr
		}
	}

	return file, nil
}

// WriteFile writes the given data the file.
//
// This function uses the [LocalPath.OpenFile] function to create/open the file before writing to it. It automatically
// closes the file after writing to it.
//
// This function may return any of the following errors:
//   - [*PathChmodError]: there was an error while changing the permissions on the file/parent folder
//   - [*PathChownError]: there was an error while changing ownership of the file/parent folder
//   - [*PathCreateError]: there was an error while creating the parent folder
//   - [*PathError]: there was a general error while working with the path
//   - [*PathOpenFileError]: there was an error while opening the file
//   - [*PathWriteError]: there was an error while writing to the file
func (p *LocalPath) WriteFile(data []byte, overwrite bool) xerrors.Error {
	flags := os.O_CREATE | os.O_RDWR
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}

	handle, xerr := p.OpenFile(flags)
	if xerr != nil {
		return xerr
	}

	defer func() {
		_ = handle.Close()
	}()

	_, err := handle.Write(data)
	if err != nil {
		return newPathWriteError(err, "failed to write to file '%s': %s", p.FSPath, err.Error()).
			WithAttr("file", p.FSPath)
	}

	return nil
}
