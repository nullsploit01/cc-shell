package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	s.handleInterrupt()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		s.cmd.OutOrStdout().Write([]byte("> "))

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				s.cmd.ErrOrStderr().Write([]byte("failed to read input: " + err.Error() + "\n"))
			}
			continue
		}

		command := scanner.Text()
		if strings.TrimSpace(command) == "" {
			continue
		}

		parts := strings.Fields(command)
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "exit":
			s.cmd.OutOrStdout().Write([]byte("Exitting.. Bye!" + "\n"))
			return nil

		case "ls":
			files, err := s.listFiles()
			if err != nil {
				s.cmd.ErrOrStderr().Write([]byte(err.Error() + "\n"))
			} else {
				for _, file := range files {
					s.cmd.OutOrStdout().Write([]byte(file + " "))
				}
				s.cmd.OutOrStdout().Write([]byte("\n"))
			}

		case "pwd":
			dir, err := s.getCurrentDir()
			if err != nil {
				s.cmd.ErrOrStderr().Write([]byte(err.Error() + "\n"))
			}

			s.cmd.OutOrStdout().Write([]byte(dir + "\n"))

		case "cd":
			if len(args) == 0 {
				s.cmd.ErrOrStderr().Write([]byte("usage: cd <directory>" + "\n"))
			} else {
				err := s.changeDirectory(args[0])
				if err != nil {
					s.cmd.ErrOrStderr().Write([]byte(err.Error() + "\n"))
				}
			}

		default:
			s.cmd.OutOrStdout().Write([]byte("no such file or directory (os error 2)\n"))
		}
	}
}

func (s *Shell) handleInterrupt() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	go func() {
		for range signalChannel {
			s.cmd.OutOrStdout().Write([]byte("\n> "))
		}
	}()
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

func (s *Shell) changeDirectory(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	return nil
}
