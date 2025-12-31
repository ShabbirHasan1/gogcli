package cmd

import (
	"context"
	"strings"
	"testing"
)

func TestDriveCommand_ValidationErrors(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com"}

	cmd := newDriveMoveCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"file1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "missing --parent") {
		t.Fatalf("expected parent error, got %v", err)
	}

	cmd = newDriveShareCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"file1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "must specify") {
		t.Fatalf("expected share validation error, got %v", err)
	}

	cmd = newDriveShareCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"file1"})
	_ = cmd.Flags().Set("anyone", "true")
	_ = cmd.Flags().Set("role", "owner")
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "invalid --role") {
		t.Fatalf("expected role error, got %v", err)
	}
}

func TestDriveDeleteUnshare_NoInput(t *testing.T) {
	flags := &rootFlags{Account: "a@b.com", NoInput: true}

	cmd := newDriveDeleteCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"file1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "refusing") {
		t.Fatalf("expected refusing error, got %v", err)
	}

	cmd = newDriveUnshareCmd(flags)
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"file1", "perm1"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "refusing") {
		t.Fatalf("expected refusing error, got %v", err)
	}
}
