package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/api/drive/v3"
)

func TestResolveDriveDownloadDestPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if _, err := resolveDriveDownloadDestPath(nil, ""); err == nil {
		t.Fatalf("expected error for nil meta")
	}
	if _, err := resolveDriveDownloadDestPath(&drive.File{Name: "x"}, ""); err == nil {
		t.Fatalf("expected error for missing id")
	}
	if _, err := resolveDriveDownloadDestPath(&drive.File{Id: "id"}, ""); err == nil {
		t.Fatalf("expected error for missing name")
	}

	meta := &drive.File{Id: "id1", Name: "../file.txt"}
	path, err := resolveDriveDownloadDestPath(meta, "")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !strings.Contains(path, "id1_file.txt") {
		t.Fatalf("unexpected path: %q", path)
	}

	dir := t.TempDir()
	path, err = resolveDriveDownloadDestPath(meta, dir)
	if err != nil {
		t.Fatalf("resolve dir: %v", err)
	}
	if !strings.HasPrefix(path, dir+string(os.PathSeparator)) {
		t.Fatalf("expected path under dir, got %q", path)
	}

	outFile := filepath.Join(t.TempDir(), "custom.bin")
	path, err = resolveDriveDownloadDestPath(meta, outFile)
	if err != nil {
		t.Fatalf("resolve file: %v", err)
	}
	if path != outFile {
		t.Fatalf("expected custom path, got %q", path)
	}
}

func TestGuessMimeType_MoreCases(t *testing.T) {
	cases := map[string]string{
		"report.pdf":  "application/pdf",
		"photo.jpg":   "image/jpeg",
		"data.csv":    "text/csv",
		"note.md":     "text/markdown",
		"unknown.bin": "application/octet-stream",
	}
	for name, want := range cases {
		if got := guessMimeType(name); got != want {
			t.Fatalf("guessMimeType(%q)=%q want %q", name, got, want)
		}
	}
}
