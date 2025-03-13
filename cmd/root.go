package cmd

import (
	"os"

	"github.com/nullsploit01/cc-shell/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ccshell",
	Short: "A minimal custom command-line shell",
	Long: `CC-Shell is a simple command-line shell built with Go. 
It supports basic shell operations such as executing commands, 
navigating directories, maintaining history, and handling keyboard inputs.

Features:
- Execute common commands (ls, pwd, cd, etc.)
- Persistent command history
- Support for 'cd -' to switch to the previous directory
- Graceful Ctrl+C handling`,
	Run: func(cmd *cobra.Command, args []string) {
		s := internal.NewShell(cmd)
		if err := s.Run(); err != nil {
			cmd.OutOrStderr().Write([]byte("Error: " + err.Error() + "\n"))
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
