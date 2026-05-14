package cmd

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDryRunE2E_MutatingCommandsSkipAuthAndAPI(t *testing.T) {
	cases := []struct {
		name string
		args []string
		op   string
	}{
		{
			name: "contacts create",
			args: []string{"contacts", "create", "--given", "Smoke", "--email", "smoke@example.com"},
			op:   "contacts.create",
		},
		{
			name: "contacts update",
			args: []string{"contacts", "update", "people/123", "--given", "Smoke"},
			op:   "contacts.update",
		},
		{
			name: "docs insert",
			args: []string{"docs", "insert", "doc123", "hello"},
			op:   "docs.insert",
		},
		{
			name: "drive move",
			args: []string{"drive", "move", "file123", "--parent", "folder123"},
			op:   "drive.move",
		},
		{
			name: "drive rename",
			args: []string{"drive", "rename", "file123", "New"},
			op:   "drive.rename",
		},
		{
			name: "meet update",
			args: []string{"meet", "update", "abc-defg-hij", "--access", "open"},
			op:   "meet.spaces.patch",
		},
		{
			name: "meet end",
			args: []string{"meet", "end", "abc-defg-hij"},
			op:   "meet.spaces.end_active_conference",
		},
		{
			name: "slides create",
			args: []string{"slides", "create", "SmokeSlides"},
			op:   "slides.create",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--json", "--dry-run", "--no-input", "--access-token", "invalid-token"}, tc.args...)
			out := captureStdout(t, func() {
				_ = captureStderr(t, func() {
					if err := Execute(args); err != nil && ExitCode(err) != 0 {
						t.Fatalf("Execute: %v", err)
					}
				})
			})

			var payload struct {
				DryRun bool   `json:"dry_run"`
				Op     string `json:"op"`
			}
			if err := json.Unmarshal([]byte(out), &payload); err != nil {
				t.Fatalf("decode dry-run output: %v\nout=%q", err, out)
			}
			if !payload.DryRun || payload.Op != tc.op {
				t.Fatalf("unexpected dry-run output: %#v", payload)
			}
		})
	}
}

func TestDryRunE2E_ContactsUpdateValidatesLocalInputs(t *testing.T) {
	tempDir := t.TempDir()
	malformed := tempDir + "/malformed.json"
	unsupported := tempDir + "/unsupported.json"
	mismatch := tempDir + "/mismatch.json"
	valid := tempDir + "/valid.json"
	for path, body := range map[string]string{
		malformed:   "{",
		unsupported: `{"notAContactField":true}`,
		mismatch:    `{"resourceName":"people/other","names":[{"givenName":"Dry"}]}`,
		valid:       `{"resourceName":"people/123","names":[{"givenName":"Dry"}]}`,
	} {
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatalf("write fixture %s: %v", path, err)
		}
	}

	invalidCases := []struct {
		name string
		args []string
	}{
		{
			name: "bad birthday",
			args: []string{"contacts", "update", "people/123", "--birthday", "nope"},
		},
		{
			name: "bad custom",
			args: []string{"contacts", "update", "people/123", "--custom", "bad"},
		},
		{
			name: "bad relation",
			args: []string{"contacts", "update", "people/123", "--relation", "bad"},
		},
		{
			name: "malformed from-file",
			args: []string{"contacts", "update", "people/123", "--from-file", malformed},
		},
		{
			name: "unsupported from-file key",
			args: []string{"contacts", "update", "people/123", "--from-file", unsupported},
		},
		{
			name: "resource mismatch from-file",
			args: []string{"contacts", "update", "people/123", "--from-file", mismatch},
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"--json", "--dry-run", "--no-input", "--access-token", "invalid-token"}, tc.args...)
			_ = captureStdout(t, func() {
				_ = captureStderr(t, func() {
					if err := Execute(args); ExitCode(err) == 0 {
						t.Fatalf("expected validation failure")
					}
				})
			})
		})
	}

	t.Run("valid from-file skips auth and API", func(t *testing.T) {
		args := []string{"--json", "--dry-run", "--no-input", "--access-token", "invalid-token", "contacts", "update", "people/123", "--from-file", valid}
		out := captureStdout(t, func() {
			_ = captureStderr(t, func() {
				if err := Execute(args); err != nil && ExitCode(err) != 0 {
					t.Fatalf("Execute: %v", err)
				}
			})
		})

		var payload struct {
			DryRun bool   `json:"dry_run"`
			Op     string `json:"op"`
		}
		if err := json.Unmarshal([]byte(out), &payload); err != nil {
			t.Fatalf("decode dry-run output: %v\nout=%q", err, out)
		}
		if !payload.DryRun || payload.Op != "contacts.update" {
			t.Fatalf("unexpected dry-run output: %#v", payload)
		}
	})
}
