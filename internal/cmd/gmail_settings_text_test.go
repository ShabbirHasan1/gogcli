package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/steipete/gogcli/internal/ui"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func TestGmailSettings_TextPaths(t *testing.T) {
	origNew := newGmailService
	t.Cleanup(func() { newGmailService = origNew })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.Contains(path, "/gmail/v1/users/me/settings/delegates") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			if strings.Contains(path, "/delegates/") {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"delegateEmail":      "d@b.com",
					"verificationStatus": "accepted",
					"delegationEnabled":  true,
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"delegates": []map[string]any{
					{"delegateEmail": "d@b.com", "verificationStatus": "accepted"},
				},
			})
			return
		case strings.Contains(path, "/gmail/v1/users/me/settings/delegates") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"delegateEmail": "d@b.com", "verificationStatus": "pending"})
			return
		case strings.Contains(path, "/gmail/v1/users/me/settings/delegates/") && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return

		case strings.Contains(path, "/gmail/v1/users/me/settings/forwardingAddresses") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			if strings.Contains(path, "/forwardingAddresses/") {
				_ = json.NewEncoder(w).Encode(map[string]any{"forwardingEmail": "f@b.com", "verificationStatus": "accepted"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"forwardingAddresses": []map[string]any{
					{"forwardingEmail": "f@b.com", "verificationStatus": "accepted"},
				},
			})
			return
		case strings.Contains(path, "/gmail/v1/users/me/settings/forwardingAddresses") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"forwardingEmail": "f@b.com", "verificationStatus": "pending"})
			return
		case strings.Contains(path, "/gmail/v1/users/me/settings/forwardingAddresses/") && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return

		case strings.Contains(path, "/gmail/v1/users/me/settings/vacation") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"enableAutoReply":       false,
				"responseSubject":       "S",
				"responseBodyHtml":      "<b>hi</b>",
				"responseBodyPlainText": "hi",
				"startTime":             "111",
				"endTime":               "222",
				"restrictToContacts":    true,
				"restrictToDomain":      true,
			})
			return
		case strings.Contains(path, "/gmail/v1/users/me/settings/vacation") && r.Method == http.MethodPut:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"enableAutoReply":    true,
				"responseSubject":    "S2",
				"startTime":          "123",
				"endTime":            "456",
				"restrictToContacts": true,
				"restrictToDomain":   false,
			})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer srv.Close()

	svc, err := gmail.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(srv.Client()),
		option.WithEndpoint(srv.URL+"/"),
	)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	newGmailService = func(context.Context, string) (*gmail.Service, error) { return svc, nil }

	flags := &rootFlags{Account: "a@b.com"}

	u, uiErr := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if uiErr != nil {
		t.Fatalf("ui.New: %v", uiErr)
	}
	ctx := ui.WithUI(context.Background(), u)

	cmd := newGmailDelegatesListCmd(flags)
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delegates list: %v", err)
	}

	cmd = newGmailDelegatesGetCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"d@b.com"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delegates get: %v", err)
	}

	cmd = newGmailDelegatesAddCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"d@b.com"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delegates add: %v", err)
	}

	cmd = newGmailDelegatesRemoveCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"d@b.com"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delegates remove: %v", err)
	}

	cmd = newGmailForwardingListCmd(flags)
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("forwarding list: %v", err)
	}

	cmd = newGmailForwardingGetCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"f@b.com"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("forwarding get: %v", err)
	}

	cmd = newGmailForwardingCreateCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"f@b.com"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("forwarding create: %v", err)
	}

	cmd = newGmailForwardingDeleteCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"f@b.com"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("forwarding delete: %v", err)
	}

	cmd = newGmailVacationGetCmd(flags)
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("vacation get: %v", err)
	}

	cmd = newGmailVacationUpdateCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{
		"--enable",
		"--subject", "S2",
		"--body", "<b>hi</b>",
		"--start", "2025-01-01T00:00:00Z",
		"--end", "2025-01-02T00:00:00Z",
		"--contacts-only",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("vacation update: %v", err)
	}
}

func TestGmailVacationUpdate_EnableDisableConflict(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}
	cmd := newGmailVacationUpdateCmd(flags)
	cmd.SetArgs([]string{"--enable", "--disable"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected conflict error")
	}
}
