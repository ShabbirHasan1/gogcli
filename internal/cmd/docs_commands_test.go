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

	"github.com/steipete/gogcli/internal/outfmt"
	"github.com/steipete/gogcli/internal/ui"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func TestDocsCreateCopyCat_JSON(t *testing.T) {
	origNew := newDriveService
	origExport := driveExportDownload
	t.Cleanup(func() {
		newDriveService = origNew
		driveExportDownload = origExport
	})

	driveExportDownload = func(context.Context, *drive.Service, string, string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("doc text")),
		}, nil
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/drive/v3") {
			path = strings.TrimPrefix(path, "/drive/v3")
		}
		switch {
		case strings.HasPrefix(path, "/files/") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":       "doc1",
				"mimeType": "application/vnd.google-apps.document",
			})
			return
		case path == "/files" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          "doc1",
				"name":        "Doc",
				"mimeType":    "application/vnd.google-apps.document",
				"webViewLink": "http://example.com/doc1",
			})
			return
		case strings.Contains(path, "/files/doc1/copy") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          "doc2",
				"name":        "Copy",
				"mimeType":    "application/vnd.google-apps.document",
				"webViewLink": "http://example.com/doc2",
			})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer srv.Close()

	svc, err := drive.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(srv.Client()),
		option.WithEndpoint(srv.URL+"/"),
	)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	newDriveService = func(context.Context, string) (*drive.Service, error) { return svc, nil }

	flags := &rootFlags{Account: "a@b.com"}
	u, uiErr := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if uiErr != nil {
		t.Fatalf("ui.New: %v", uiErr)
	}
	ctx := ui.WithUI(context.Background(), u)
	ctx = outfmt.WithMode(ctx, outfmt.Mode{JSON: true})

	_ = captureStdout(t, func() {
		cmd := newDocsCreateCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"Doc"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("create: %v", err)
		}
	})

	_ = captureStdout(t, func() {
		cmd := newDocsCopyCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"doc1", "Copy"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("copy: %v", err)
		}
	})

	out := captureStdout(t, func() {
		cmd := newDocsCatCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"doc1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("cat: %v", err)
		}
	})
	if !strings.Contains(out, "doc text") {
		t.Fatalf("unexpected cat output: %q", out)
	}
}

func TestDocsCreateCopyCat_Text(t *testing.T) {
	origNew := newDriveService
	origExport := driveExportDownload
	t.Cleanup(func() {
		newDriveService = origNew
		driveExportDownload = origExport
	})

	driveExportDownload = func(context.Context, *drive.Service, string, string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("doc text")),
		}, nil
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/drive/v3") {
			path = strings.TrimPrefix(path, "/drive/v3")
		}
		switch {
		case strings.HasPrefix(path, "/files/") && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":       "doc1",
				"mimeType": "application/vnd.google-apps.document",
			})
			return
		case path == "/files" && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          "doc1",
				"name":        "Doc",
				"mimeType":    "application/vnd.google-apps.document",
				"webViewLink": "http://example.com/doc1",
			})
			return
		case strings.Contains(path, "/files/doc1/copy") && r.Method == http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":          "doc2",
				"name":        "Copy",
				"mimeType":    "application/vnd.google-apps.document",
				"webViewLink": "http://example.com/doc2",
			})
			return
		default:
			http.NotFound(w, r)
			return
		}
	}))
	defer srv.Close()

	svc, err := drive.NewService(context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(srv.Client()),
		option.WithEndpoint(srv.URL+"/"),
	)
	if err != nil {
		t.Fatalf("NewService: %v", err)
	}
	newDriveService = func(context.Context, string) (*drive.Service, error) { return svc, nil }

	flags := &rootFlags{Account: "a@b.com"}

	out := captureStdout(t, func() {
		u, uiErr := ui.New(ui.Options{Stdout: os.Stdout, Stderr: io.Discard, Color: "never"})
		if uiErr != nil {
			t.Fatalf("ui.New: %v", uiErr)
		}
		ctx := ui.WithUI(context.Background(), u)

		cmd := newDocsCreateCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"Doc"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("create: %v", err)
		}

		cmd = newDocsCopyCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"doc1", "Copy"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("copy: %v", err)
		}

		cmd = newDocsCatCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"doc1"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("cat: %v", err)
		}
	})
	if !strings.Contains(out, "doc text") || !strings.Contains(out, "id\tdoc1") {
		t.Fatalf("unexpected output: %q", out)
	}
}
