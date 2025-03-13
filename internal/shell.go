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

		if command != "history" {
			s.saveToHistory(command)
		}

		parts := strings.Fields(command)
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "exit":
			s.cmd.OutOrStdout().Write([]byte("Exitting.. Bye!" + "\n"))
			return nil

		case "ls":
			s.printFilesInCurrentDirectory()

		case "pwd":
			s.printWorkingDirectory()

		case "cd":
			s.changeDirectory(args)

		case "history":
			s.printHistory()

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

func (s *Shell) saveToHistory(command string) {
	s.history = append(s.history, command)

	file, err := os.OpenFile(s.histPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.cmd.ErrOrStderr().Write([]byte("Failed to save history: " + err.Error() + "\n"))
		return
	}
	defer file.Close()
	file.WriteString(command + "\n")
}

func (s *Shell) printHistory() {
	for i, cmd := range s.history {
		s.cmd.OutOrStdout().Write(fmt.Appendf(nil, "%d %s\n", i+1, cmd))
	}
}

func (s *Shell) printFilesInCurrentDirectory() {
	files, err := s.listFiles()
	if err != nil {
		s.cmd.ErrOrStderr().Write([]byte(err.Error() + "\n"))
	} else {
		for _, file := range files {
			s.cmd.OutOrStdout().Write([]byte(file + " "))
		}
		s.cmd.OutOrStdout().Write([]byte("\n"))
	}
}

func (s *Shell) printWorkingDirectory() {
	dir, err := s.getCurrentDir()
	if err != nil {
		s.cmd.ErrOrStderr().Write([]byte(err.Error() + "\n"))
	}

	s.cmd.OutOrStdout().Write([]byte(dir + "\n"))
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
