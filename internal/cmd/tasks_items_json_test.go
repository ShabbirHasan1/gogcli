package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/steipete/gogcli/internal/outfmt"
	"github.com/steipete/gogcli/internal/ui"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func TestTasksItems_JSONPaths(t *testing.T) {
	origNew := newTasksService
	t.Cleanup(func() { newTasksService = origNew })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/tasks/v1/lists/l1/tasks") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"items":         []map[string]any{{"id": "t1", "title": "Task"}},
				"nextPageToken": "next",
			})
			return
		case strings.HasSuffix(r.URL.Path, "/tasks/v1/lists/l1/tasks") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "t1", "title": "Task"})
			return
		case strings.Contains(r.URL.Path, "/tasks/v1/lists/l1/tasks/t1") && r.Method == http.MethodPatch:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "t1", "title": "Task", "status": "completed"})
			return
		case strings.Contains(r.URL.Path, "/tasks/v1/lists/l1/tasks/t1") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "t1", "title": "Task"})
			return
		case strings.Contains(r.URL.Path, "/tasks/v1/lists/l1/tasks/t1") && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			return
		case r.URL.Path == "/tasks/v1/lists/l1/clear" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusNoContent)
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer srv.Close()

	svc, err := tasks.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(srv.Client()),
		option.WithEndpoint(srv.URL+"/"),
	)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	newTasksService = func(context.Context, string) (*tasks.Service, error) { return svc, nil }

	flags := &rootFlags{Account: "a@b.com", Force: true}
	u, uiErr := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if uiErr != nil {
		t.Fatalf("ui.New: %v", uiErr)
	}
	ctx := ui.WithUI(context.Background(), u)
	ctx = outfmt.WithMode(ctx, outfmt.Mode{JSON: true})

	// list
	_ = captureStdout(t, func() {
		cmd := newTasksListCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1"})
		_ = cmd.Flags().Set("due-min", "2025-01-01T00:00:00Z")
		_ = cmd.Flags().Set("due-max", "2025-01-02T00:00:00Z")
		_ = cmd.Flags().Set("completed-min", "2025-01-01T00:00:00Z")
		_ = cmd.Flags().Set("completed-max", "2025-01-02T00:00:00Z")
		_ = cmd.Flags().Set("updated-min", "2025-01-01T00:00:00Z")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("list: %v", err)
		}
	})

	// add
	_ = captureStdout(t, func() {
		cmd := newTasksAddCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1"})
		_ = cmd.Flags().Set("title", "Task")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("add: %v", err)
		}
	})

	// update
	_ = captureStdout(t, func() {
		cmd := newTasksUpdateCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1", "t1"})
		_ = cmd.Flags().Set("status", "completed")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("update: %v", err)
		}
	})

	// done
	_ = captureStdout(t, func() {
		cmd := newTasksDoneCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1", "t1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("done: %v", err)
		}
	})

	// undo
	_ = captureStdout(t, func() {
		cmd := newTasksUndoCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1", "t1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("undo: %v", err)
		}
	})

	// delete
	_ = captureStdout(t, func() {
		cmd := newTasksDeleteCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1", "t1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("delete: %v", err)
		}
	})

	// clear
	_ = captureStdout(t, func() {
		cmd := newTasksClearCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"l1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("clear: %v", err)
		}
	})
}

func TestTasksAddCmd_MissingTitle(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}
	cmd := newTasksAddCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"l1"})
	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error")
	}
}
