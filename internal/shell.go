package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

const historyFile = ".ccshell_history"

type Shell struct {
	cmd      *cobra.Command
	history  []string
	histPath string
	prevDir  string
}

func NewShell(cmd *cobra.Command) *Shell {
	homeDir, _ := os.UserHomeDir()
	histPath := filepath.Join(homeDir, historyFile)

	shell := &Shell{
		cmd:      cmd,
		histPath: histPath,
	}

	shell.loadHistory()
	return shell
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
			s.changeDirectory(args)

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

func (s *Shell) loadHistory() {
	file, err := os.Open(s.histPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s.history = append(s.history, scanner.Text())
	}
}

func (s *Shell) changeDirectory(args []string) {
	var targetDir string

	if len(args) == 0 || args[0] == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			s.cmd.ErrOrStderr().Write([]byte("Failed to get home directory: " + err.Error() + "\n"))
			return
		}
		targetDir = homeDir
	} else if args[0] == "-" {
		if s.prevDir == "" {
			s.cmd.OutOrStdout().Write([]byte("OLDPWD not set\n"))
			return
		}
		targetDir = s.prevDir
	} else {
		targetDir = args[0]
	}

	currDir, err := os.Getwd()
	if err != nil {
		s.cmd.ErrOrStderr().Write([]byte("Failed to get current directory: " + err.Error() + "\n"))
		return
	}

	err = os.Chdir(targetDir)
	if err != nil {
		s.cmd.ErrOrStderr().Write([]byte("Failed to change directory: " + err.Error() + "\n"))
		return
	}

	s.prevDir = currDir
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
