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
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func TestTasks_TextPaths(t *testing.T) {
	origNew := newTasksService
	t.Cleanup(func() { newTasksService = origNew })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/tasks/v1/users/@me/lists" && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{{"id": "l1", "title": "List"}},
			})
			return
		case r.URL.Path == "/tasks/v1/users/@me/lists" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"id": "l1", "title": "List"})
			return
		case strings.HasSuffix(r.URL.Path, "/tasks/v1/lists/l1/tasks") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{
					{"id": "t1", "title": "Task", "status": "needsAction", "due": "2025-01-01T00:00:00Z", "updated": "2025-01-01T00:00:00Z"},
				},
			})
			return
		case strings.HasSuffix(r.URL.Path, "/tasks/v1/lists/l1/tasks") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          "t1",
				"title":       "Task",
				"status":      "needsAction",
				"due":         "2025-01-01T00:00:00Z",
				"webViewLink": "http://example.com/task",
			})
			return
		case strings.Contains(r.URL.Path, "/tasks/v1/lists/l1/tasks/t1") && r.Method == http.MethodPatch:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          "t1",
				"title":       "Task",
				"status":      "completed",
				"due":         "2025-01-01T00:00:00Z",
				"webViewLink": "http://example.com/task",
			})
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

	cmd := newTasksListsCmd(flags)
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lists: %v", err)
	}

	cmd = newTasksListsCreateCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"My", "List"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lists create: %v", err)
	}

	cmd = newTasksListCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}

	cmd = newTasksAddCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1"})
	cmd.Flags().Set("title", "Task")
	cmd.Flags().Set("notes", "Notes")
	cmd.Flags().Set("due", "2025-01-01T00:00:00Z")
	cmd.Flags().Set("parent", "p1")
	cmd.Flags().Set("previous", "p0")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("add: %v", err)
	}

	cmd = newTasksUpdateCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1", "t1"})
	cmd.Flags().Set("title", "New title")
	cmd.Flags().Set("notes", "New notes")
	cmd.Flags().Set("due", "2025-01-02T00:00:00Z")
	cmd.Flags().Set("status", "completed")
	if err := cmd.Execute(); err != nil {
		t.Fatalf("update: %v", err)
	}

	cmd = newTasksDoneCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1", "t1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("done: %v", err)
	}

	cmd = newTasksUndoCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1", "t1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("undo: %v", err)
	}

	cmd = newTasksDeleteCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1", "t1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delete: %v", err)
	}

	cmd = newTasksClearCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("clear: %v", err)
	}
}

func TestTasksLists_NoItems(t *testing.T) {
	origNew := newTasksService
	t.Cleanup(func() { newTasksService = origNew })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/tasks/v1/users/@me/lists" && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"items": []map[string]any{}})
			return
		case strings.HasSuffix(r.URL.Path, "/tasks/v1/lists/l1/tasks") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"items": []map[string]any{}})
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

	flags := &rootFlags{Account: "a@b.com"}
	u, uiErr := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if uiErr != nil {
		t.Fatalf("ui.New: %v", uiErr)
	}
	ctx := ui.WithUI(context.Background(), u)

	cmd := newTasksListsCmd(flags)
	cmd.SetContext(ctx)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("lists: %v", err)
	}

	cmd = newTasksListCmd(flags)
	cmd.SetContext(ctx)
	cmd.SetArgs([]string{"l1"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
}
