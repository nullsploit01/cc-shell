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

		switch command {
		case "exit":
			fmt.Println("Exitting.. Bye!")
			return nil

		case "ls":
			files, err := s.listFiles()
			if err != nil {
				return err
			} else {
				for _, file := range files {
					s.cmd.OutOrStdout().Write([]byte(file + " "))
				}
				s.cmd.OutOrStdout().Write([]byte("\n"))
			}

		case "pwd":
			dir, err := s.getCurrentDir()
			if err != nil {
				return err
			}

			s.cmd.OutOrStdout().Write([]byte(dir + "\n"))

		default:
			return fmt.Errorf("no such file or directory (os error 2)")
		}
	}

	return nil
}

func (s *Shell) listFiles() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range files {
		name := file.Name()
		if file.IsDir() {
			name += "/" // Append '/' for directories
		}
		fileList = append(fileList, name)
	}

	return fileList, nil
}

func (s *Shell) getCurrentDir() (string, error) {
	return os.Getwd()
}
