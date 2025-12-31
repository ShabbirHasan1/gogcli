package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCompletionCmd(t *testing.T) {
	cases := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range cases {
		shell := shell
		t.Run(shell, func(t *testing.T) {
			root := &cobra.Command{Use: "gog"}
			root.AddCommand(newCompletionCmd())

			out := captureStdout(t, func() {
				root.SetArgs([]string{"completion", shell})
				if err := root.Execute(); err != nil {
					t.Fatalf("execute: %v", err)
				}
			})
			if out == "" {
				t.Fatalf("expected completion output")
			}
		})
	}
}
