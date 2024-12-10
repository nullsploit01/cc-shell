package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Shell struct{}

func NewShell() *Shell {
	return &Shell{}
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
		fmt.Println(command)
	}

	return nil
}
