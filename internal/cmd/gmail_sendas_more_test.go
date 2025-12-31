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
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func TestGmailSendAs_VerifyDeleteUpdate_JSON(t *testing.T) {
	origNew := newGmailService
	t.Cleanup(func() { newGmailService = origNew })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/settings/sendAs/") && strings.HasSuffix(r.URL.Path, "/verify") && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if strings.Contains(r.URL.Path, "/settings/sendAs/") && r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if strings.Contains(r.URL.Path, "/settings/sendAs/") && !strings.HasSuffix(r.URL.Path, "/verify") {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"sendAsEmail":        "work@company.com",
				"displayName":        "Work Alias",
				"replyToAddress":     "reply@company.com",
				"signature":          "Sig",
				"treatAsAlias":       true,
				"verificationStatus": "accepted",
			})
			return
		}
		http.NotFound(w, r)
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
	ctx = outfmt.WithMode(ctx, outfmt.Mode{JSON: true})

	// verify
	_ = captureStdout(t, func() {
		cmd := newGmailSendAsVerifyCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"work@company.com"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("verify: %v", err)
		}
	})

	// update
	_ = captureStdout(t, func() {
		cmd := newGmailSendAsUpdateCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"work@company.com"})
		_ = cmd.Flags().Set("display-name", "Work Alias")
		_ = cmd.Flags().Set("reply-to", "reply@company.com")
		_ = cmd.Flags().Set("signature", "Sig")
		_ = cmd.Flags().Set("treat-as-alias", "true")
		if err := cmd.Execute(); err != nil {
			t.Fatalf("update: %v", err)
		}
	})

	// delete
	_ = captureStdout(t, func() {
		cmd := newGmailSendAsDeleteCmd(flags)
		cmd.SetContext(ctx)
		cmd.SetArgs([]string{"work@company.com"})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("delete: %v", err)
		}
	})
}
