package cmd

import (
	"context"
	"testing"
)

func TestGmailWatchServeCmd_ValidationErrors(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}

	t.Run("path missing slash", func(t *testing.T) {
		cmd := newGmailWatchServeCmd(flags)
		cmd.SetContext(context.Background())
		cmd.SetArgs([]string{"--path", "hook", "--port", "9999"})
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("port missing", func(t *testing.T) {
		cmd := newGmailWatchServeCmd(flags)
		cmd.SetContext(context.Background())
		cmd.SetArgs([]string{"--port", "0"})
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("non-loopback without auth", func(t *testing.T) {
		cmd := newGmailWatchServeCmd(flags)
		cmd.SetContext(context.Background())
		cmd.SetArgs([]string{"--bind", "0.0.0.0", "--port", "9999"})
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("oidc email requires verify", func(t *testing.T) {
		cmd := newGmailWatchServeCmd(flags)
		cmd.SetContext(context.Background())
		cmd.SetArgs([]string{"--oidc-email", "svc@example.com", "--port", "9999"})
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("oidc audience requires verify", func(t *testing.T) {
		cmd := newGmailWatchServeCmd(flags)
		cmd.SetContext(context.Background())
		cmd.SetArgs([]string{"--oidc-audience", "aud", "--port", "9999"})
		if err := cmd.Execute(); err == nil {
			t.Fatalf("expected error")
		}
	})
}
