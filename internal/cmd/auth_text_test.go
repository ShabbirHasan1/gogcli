package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/steipete/gogcli/internal/secrets"
)

func TestAuthTextOutputs(t *testing.T) {
	origOpen := openSecretsStore
	t.Cleanup(func() { openSecretsStore = origOpen })

	store := newMemSecretsStore()
	openSecretsStore = func() (secrets.Store, error) { return store, nil }

	if err := store.SetToken("a@b.com", secrets.Token{
		Services:     []string{"gmail"},
		RefreshToken: "rt",
		CreatedAt:    time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
	}); err != nil {
		t.Fatalf("SetToken: %v", err)
	}

	listOut := captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"auth", "list"}); err != nil {
				t.Fatalf("list: %v", err)
			}
		})
	})
	if !strings.Contains(listOut, "a@b.com") {
		t.Fatalf("unexpected list output: %q", listOut)
	}

	keysOut := captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"auth", "tokens", "list"}); err != nil {
				t.Fatalf("tokens list: %v", err)
			}
		})
	})
	if !strings.Contains(keysOut, "token:a@b.com") {
		t.Fatalf("unexpected keys output: %q", keysOut)
	}

	delOut := captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"--force", "auth", "tokens", "delete", "a@b.com"}); err != nil {
				t.Fatalf("tokens delete: %v", err)
			}
		})
	})
	if !strings.Contains(delOut, "deleted") {
		t.Fatalf("unexpected delete output: %q", delOut)
	}

	// Re-add and remove via auth remove.
	if err := store.SetToken("a@b.com", secrets.Token{RefreshToken: "rt"}); err != nil {
		t.Fatalf("SetToken: %v", err)
	}
	rmOut := captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"--force", "auth", "remove", "a@b.com"}); err != nil {
				t.Fatalf("remove: %v", err)
			}
		})
	})
	if !strings.Contains(rmOut, "deleted") {
		t.Fatalf("unexpected remove output: %q", rmOut)
	}
}

func TestAuthListAndTokens_NoTokens_Text(t *testing.T) {
	origOpen := openSecretsStore
	t.Cleanup(func() { openSecretsStore = origOpen })

	store := newMemSecretsStore()
	openSecretsStore = func() (secrets.Store, error) { return store, nil }

	errOut := captureStderr(t, func() {
		_ = captureStdout(t, func() {
			if err := Execute([]string{"auth", "list"}); err != nil {
				t.Fatalf("list: %v", err)
			}
		})
	})
	if !strings.Contains(errOut, "No tokens stored") {
		t.Fatalf("unexpected stderr: %q", errOut)
	}

	errOut = captureStderr(t, func() {
		_ = captureStdout(t, func() {
			if err := Execute([]string{"auth", "tokens", "list"}); err != nil {
				t.Fatalf("tokens list: %v", err)
			}
		})
	})
	if !strings.Contains(errOut, "No tokens stored") {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestAuthTokensExportImport_Text(t *testing.T) {
	origOpen := openSecretsStore
	t.Cleanup(func() { openSecretsStore = origOpen })

	store := newMemSecretsStore()
	openSecretsStore = func() (secrets.Store, error) { return store, nil }

	if err := store.SetToken("a@b.com", secrets.Token{
		RefreshToken: "rt",
	}); err != nil {
		t.Fatalf("SetToken: %v", err)
	}

	outPath := filepath.Join(t.TempDir(), "token.json")
	out := captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"auth", "tokens", "export", "a@b.com", "--out", outPath}); err != nil {
				t.Fatalf("export: %v", err)
			}
		})
	})
	if !strings.Contains(out, "exported") || !strings.Contains(out, outPath) {
		t.Fatalf("unexpected export output: %q", out)
	}

	if err := store.DeleteToken("a@b.com"); err != nil {
		t.Fatalf("DeleteToken: %v", err)
	}

	out = captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"auth", "tokens", "import", outPath}); err != nil {
				t.Fatalf("import: %v", err)
			}
		})
	})
	if !strings.Contains(out, "imported") {
		t.Fatalf("unexpected import output: %q", out)
	}
}

func TestAuthCredentials_Text(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	in := filepath.Join(t.TempDir(), "creds.json")
	if err := os.WriteFile(in, []byte(`{"installed":{"client_id":"id","client_secret":"sec"}}`), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	out := captureStdout(t, func() {
		_ = captureStderr(t, func() {
			if err := Execute([]string{"auth", "credentials", in}); err != nil {
				t.Fatalf("Execute: %v", err)
			}
		})
	})
	if !strings.Contains(out, "path") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestAuthAdd_InvalidService(t *testing.T) {
	if err := Execute([]string{"auth", "add", "a@b.com", "--services", "nope"}); err == nil {
		t.Fatalf("expected error")
	}
}
