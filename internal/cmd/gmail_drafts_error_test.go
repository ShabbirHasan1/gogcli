package cmd

import (
	"context"
	"strings"
	"testing"
)

func TestGmailDraftsCreate_ValidationErrors(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}

	cmd := newGmailDraftsCreateCmd(flags)
	cmd.SetContext(context.Background())
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "required: --to, --subject") {
		t.Fatalf("expected required to/subject error, got %v", err)
	}

	cmd = newGmailDraftsCreateCmd(flags)
	cmd.SetContext(context.Background())
	_ = cmd.Flags().Set("to", "b@b.com")
	_ = cmd.Flags().Set("subject", "Hi")
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "required: --body") {
		t.Fatalf("expected body error, got %v", err)
	}
}
