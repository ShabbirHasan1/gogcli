package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfirmDestructive_NoInput(t *testing.T) {
	flags := &rootFlags{NoInput: true}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	err := confirmDestructive(cmd, flags, "delete something")
	if err == nil || !strings.Contains(err.Error(), "refusing") {
		t.Fatalf("expected refusing error, got %v", err)
	}
}

func TestConfirmDestructive_Force(t *testing.T) {
	flags := &rootFlags{Force: true}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	if err := confirmDestructive(cmd, flags, "delete something"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
