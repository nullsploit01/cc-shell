package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Shell struct {
	cmd *cobra.Command
}

func NewShell(cmd *cobra.Command) *Shell {
	return &Shell{
		cmd: cmd,
	}
}

func (s *Shell) Run() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			break
		}

		command := scanner.Text()
		if strings.TrimSpace(command) == "" {
			continue
		}

		if strings.Compare(command, "exit") == 0 {
			s.cmd.Println("Exitting.. Bye!")
			break
		}

		return fmt.Errorf("No such file or directory (os error 2)")
	}

	return nil
}
