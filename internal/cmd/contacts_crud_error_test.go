package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/steipete/gogcli/internal/ui"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func TestContactsListAndGet_NoResults_Text(t *testing.T) {
	origContacts := newPeopleContactsService
	t.Cleanup(func() { newPeopleContactsService = origContacts })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "people/me/connections") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"connections": []map[string]any{}})
			return
		case strings.Contains(r.URL.Path, "people:searchContacts") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"results": []map[string]any{}})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer srv.Close()

	svc, err := people.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(srv.Client()),
		option.WithEndpoint(srv.URL+"/"),
	)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	newPeopleContactsService = func(context.Context, string) (*people.Service, error) { return svc, nil }

	flags := &rootFlags{Account: "a@b.com"}
	errOut := captureStderr(t, func() {
		_ = captureStdout(t, func() {
			u, uiErr := ui.New(ui.Options{Stdout: io.Discard, Stderr: os.Stderr, Color: "never"})
			if uiErr != nil {
				t.Fatalf("ui.New: %v", uiErr)
			}
			ctx := ui.WithUI(context.Background(), u)

			cmd := newContactsListCmd(flags)
			cmd.SetContext(ctx)
			if err := cmd.Execute(); err != nil {
				t.Fatalf("list: %v", err)
			}

			cmd = newContactsGetCmd(flags)
			cmd.SetContext(ctx)
			cmd.SetArgs([]string{"missing@example.com"})
			if err := cmd.Execute(); err != nil {
				t.Fatalf("get: %v", err)
			}
		})
	})
	if !strings.Contains(errOut, "No contacts") && !strings.Contains(errOut, "Not found") {
		t.Fatalf("unexpected stderr: %q", errOut)
	}
}

func TestContactsUpdateDelete_InvalidResource(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}

	cmd := newContactsUpdateCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"nope"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "resourceName must start") {
		t.Fatalf("expected resourceName error, got %v", err)
	}

	cmd = newContactsDeleteCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"nope"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "resourceName must start") {
		t.Fatalf("expected resourceName error, got %v", err)
	}
}
