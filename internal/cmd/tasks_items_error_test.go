package cmd

import (
	"context"
	"strings"
	"testing"
)

func TestTasksUpdate_ValidationErrors(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}

	cmd := newTasksUpdateCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"l1", "t1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "no fields to update") {
		t.Fatalf("expected no fields error, got %v", err)
	}

	cmd = newTasksUpdateCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"l1", "t1"})
	_ = cmd.Flags().Set("status", "nope")
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "invalid --status") {
		t.Fatalf("expected status error, got %v", err)
	}
}
