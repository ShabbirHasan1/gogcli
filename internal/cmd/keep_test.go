package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	keepapi "google.golang.org/api/keep/v1"

	"github.com/steipete/gogcli/internal/config"
)

func TestGetKeepService_NoServiceAccountConfigured(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg-config"))

	orig := newKeepServiceWithSA
	t.Cleanup(func() { newKeepServiceWithSA = orig })

	called := false
	newKeepServiceWithSA = func(context.Context, string, string) (*keepapi.Service, error) {
		called = true
		return &keepapi.Service{}, nil
	}

	_, err := getKeepService(context.Background(), &RootFlags{Account: "a@b.com"}, &KeepCmd{})
	if err == nil {
		t.Fatalf("expected error")
	}
	if ExitCode(err) != 2 {
		t.Fatalf("expected exit code 2, got %v", ExitCode(err))
	}
	if called {
		t.Fatalf("expected no service account usage")
	}
}

func TestGetKeepService_UsesStoredServiceAccount(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, "xdg-config"))

	account := "a@b.com"
	saPath, err := config.KeepServiceAccountPath(account)
	if err != nil {
		t.Fatalf("KeepServiceAccountPath: %v", err)
	}
	if mkdirErr := os.MkdirAll(filepath.Dir(saPath), 0o700); mkdirErr != nil {
		t.Fatalf("mkdir: %v", mkdirErr)
	}
	if writeErr := os.WriteFile(saPath, []byte("{}"), 0o600); writeErr != nil {
		t.Fatalf("write: %v", writeErr)
	}

	orig := newKeepServiceWithSA
	t.Cleanup(func() { newKeepServiceWithSA = orig })

	var gotPath, gotImpersonate string
	newKeepServiceWithSA = func(ctx context.Context, path, impersonate string) (*keepapi.Service, error) {
		gotPath = path
		gotImpersonate = impersonate
		return &keepapi.Service{}, nil
	}

	svc, err := getKeepService(context.Background(), &RootFlags{Account: account}, &KeepCmd{})
	if err != nil {
		t.Fatalf("getKeepService: %v", err)
	}
	if svc == nil {
		t.Fatalf("expected service")
	}
	if gotPath != saPath {
		t.Fatalf("unexpected path: %q", gotPath)
	}
	if gotImpersonate != account {
		t.Fatalf("unexpected impersonate: %q", gotImpersonate)
	}
}
