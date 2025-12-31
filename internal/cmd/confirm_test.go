package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfirmDestructive_Force(t *testing.T) {
	cmd := &cobra.Command{Use: "cmd"}
	cmd.SetContext(context.Background())
	flags := &rootFlags{Force: true}

	if err := confirmDestructive(cmd, flags, "delete"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestConfirmDestructive_NoInput(t *testing.T) {
	cmd := &cobra.Command{Use: "cmd"}
	cmd.SetContext(context.Background())
	flags := &rootFlags{NoInput: true}

	err := confirmDestructive(cmd, flags, "delete")
	if err == nil {
		t.Fatalf("expected error")
	}
	if ExitCode(err) != 2 {
		t.Fatalf("expected exit code 2, got %d", ExitCode(err))
	}
	if !strings.Contains(err.Error(), "refusing to delete") {
		t.Fatalf("unexpected error: %v", err)
	}
}
