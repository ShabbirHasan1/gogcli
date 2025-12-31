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
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func TestCalendarMoreCommands_Text(t *testing.T) {
	origNew := newCalendarService
	t.Cleanup(func() { newCalendarService = origNew })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/calendar/v3") {
			path = strings.TrimPrefix(path, "/calendar/v3")
		}
		switch {
		case path == "/users/me/calendarList" && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{
					{"id": "cal1", "summary": "Primary", "accessRole": "owner"},
				},
			})
			return
		case strings.Contains(path, "/calendars/cal1/acl") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{
					{"role": "reader", "scope": map[string]any{"type": "user", "value": "a@b.com"}},
				},
			})
			return
		case strings.Contains(path, "/calendars/cal1/events/evt1") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":       "evt1",
				"summary":  "Old",
				"start":    map[string]any{"dateTime": "2025-01-01T10:00:00Z"},
				"end":      map[string]any{"dateTime": "2025-01-01T11:00:00Z"},
				"htmlLink": "http://example.com/event",
			})
			return
		case strings.Contains(path, "/calendars/cal1/events/evt1") && (r.Method == http.MethodPut || r.Method == http.MethodPatch):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":       "evt1",
				"summary":  "Updated",
				"htmlLink": "http://example.com/event",
			})
			return
		case strings.Contains(path, "/calendars/cal1/events/evt1") && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		case strings.Contains(path, "/calendars/cal1/events") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":       "evt2",
				"summary":  "Created",
				"htmlLink": "http://example.com/created",
			})
			return
		case strings.Contains(strings.ToLower(path), "freebusy") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"calendars": map[string]any{
					"cal1": map[string]any{
						"busy": []map[string]any{{"start": "2025-01-01T10:00:00Z", "end": "2025-01-01T11:00:00Z"}},
					},
				},
			})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer srv.Close()

	svc, err := calendar.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(srv.Client()),
		option.WithEndpoint(srv.URL+"/"),
	)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	newCalendarService = func(context.Context, string) (*calendar.Service, error) { return svc, nil }

	flags := &rootFlags{Account: "a@b.com", Force: true}

	out := captureStdout(t, func() {
		u, uiErr := ui.New(ui.Options{Stdout: os.Stdout, Stderr: io.Discard, Color: "never"})
		if uiErr != nil {
			t.Fatalf("ui.New: %v", uiErr)
		}
		ctx := ui.WithUI(context.Background(), u)

		cmd := newCalendarCalendarsCmd(flags)
		cmd.SetContext(ctx)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("calendars: %v", err)
		}

		cmd = newCalendarAclCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"cal1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("acl: %v", err)
		}

		cmd = newCalendarEventCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"cal1", "evt1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("event: %v", err)
		}

		cmd = newCalendarCreateCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"cal1"})
		_ = cmd.Flags().Set("summary", "Created")
		_ = cmd.Flags().Set("from", "2025-01-01T12:00:00Z")
		_ = cmd.Flags().Set("to", "2025-01-01T13:00:00Z")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("create: %v", err)
		}

		cmd = newCalendarUpdateCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"cal1", "evt1"})
		_ = cmd.Flags().Set("summary", "Updated")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("update: %v", err)
		}

		cmd = newCalendarDeleteCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"cal1", "evt1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("delete: %v", err)
		}

		cmd = newCalendarFreeBusyCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"cal1"})
		_ = cmd.Flags().Set("from", "2025-01-01T00:00:00Z")
		_ = cmd.Flags().Set("to", "2025-01-02T00:00:00Z")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("freebusy: %v", err)
		}
	})
	if !strings.Contains(out, "CALENDAR") || !strings.Contains(out, "evt1") {
		t.Fatalf("unexpected output: %q", out)
	}
}
