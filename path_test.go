package xtypes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalPathHappyPath(t *testing.T) {
	tmp := t.TempDir()
	dirTarget := filepath.Join(tmp, "a")
	target := filepath.Join(tmp, "a", "b.txt")

	dirPath := &LocalPath{
		AutoChmod: true,
		AutoChown: true,
		DirMode:   FileMode(0o755),
		FSPath:    dirTarget,
		Group:     GroupID(os.Getgid()),
		Owner:     UserID(os.Getuid()),
	}

	lp := &LocalPath{
		AutoCreateParent: true,
		AutoChmod:        true,
		AutoChown:        true,
		DirMode:          FileMode(0o755),
		FileMode:         FileMode(0o644),
		FSPath:           target,
		Group:            GroupID(os.Getgid()),
		Owner:            UserID(os.Getuid()),
	}

	if xerr := lp.ToAbs(); xerr != nil {
		t.Fatalf("ToAbs: %v", xerr)
	}
	if !filepath.IsAbs(lp.FSPath) {
		t.Fatalf("expected absolute path")
	}
	if len(lp.Attrs()) == 0 {
		t.Fatalf("attrs empty")
	}

	if xerr := dirPath.MkdirAll(); xerr != nil {
		t.Fatalf("MkdirAll: %v", xerr)
	}

	f, xerr := lp.OpenFile(os.O_CREATE | os.O_RDWR)
	if xerr != nil {
		t.Fatalf("OpenFile: %v", xerr)
	}
	_ = f.Close()

	if xerr := lp.WriteFile([]byte("hello"), true); xerr != nil {
		t.Fatalf("WriteFile overwrite: %v", xerr)
	}
	if xerr := lp.WriteFile([]byte(" world"), false); xerr != nil {
		t.Fatalf("WriteFile append: %v", xerr)
	}

	content, err := os.ReadFile(lp.FSPath)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(content) != "hello world" {
		t.Fatalf("unexpected content: %q", string(content))
	}

	if xerr := lp.Chmod(); xerr != nil {
		t.Fatalf("Chmod: %v", xerr)
	}
	if xerr := lp.Chown(); xerr != nil {
		t.Fatalf("Chown: %v", xerr)
	}
}

func TestLocalPathErrorPaths(t *testing.T) {
	bad := &LocalPath{FSPath: string([]byte{'b', 'a', 'd', 0})}
	if xerr := bad.MkdirAll(); xerr == nil {
		t.Fatalf("expected MkdirAll error")
	}
	if _, xerr := bad.OpenFile(os.O_CREATE | os.O_RDWR); xerr == nil {
		t.Fatalf("expected OpenFile error")
	}
	if xerr := bad.WriteFile([]byte("x"), true); xerr == nil {
		t.Fatalf("expected WriteFile error")
	}

	missing := &LocalPath{FSPath: filepath.Join(t.TempDir(), "missing")}
	xerr := missing.Chmod()
	if xerr == nil {
		t.Fatalf("expected Chmod PathError")
	}
	if !strings.Contains(xerr.Error(), "failed to change permissions") {
		t.Fatalf("unexpected chmod error: %v", xerr)
	}
}
