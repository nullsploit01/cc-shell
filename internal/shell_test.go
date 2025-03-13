package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func newTestShell() *Shell {
	cmd := &cobra.Command{}
	return NewShell(cmd)
}

func TestSaveToHistory(t *testing.T) {
	shell := newTestShell()

	tmpFile, err := os.CreateTemp("", "test_history")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	shell.histPath = tmpFile.Name()

	shell.saveToHistory("echo hello")

	data, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)

	assert.Contains(t, string(data), "echo hello")
}

func TestLoadHistoryContains(t *testing.T) {
	shell := newTestShell()

	tmpFile, err := os.CreateTemp("", "test_history")
	assert.NoError(t, err)

	historyContent := "ls\npwd\necho hello\necho world\ncd Documents\n"
	err = os.WriteFile(tmpFile.Name(), []byte(historyContent), 0644)
	assert.NoError(t, err)

	shell.histPath = tmpFile.Name()

	shell.loadHistory()

	assert.Contains(t, shell.history, "echo hello")
	assert.Contains(t, shell.history, "echo world")
}

func TestChangeDirectory(t *testing.T) {
	shell := newTestShell()

	initialDir, err := os.Getwd()
	assert.NoError(t, err)

	shell.changeDirectory([]string{"~"})
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)
	currDir, err := os.Getwd()
	assert.NoError(t, err)
	assert.Equal(t, homeDir, currDir)

	shell.changeDirectory([]string{initialDir})
	currDir, err = os.Getwd()
	assert.NoError(t, err)
	assert.Equal(t, initialDir, currDir)
}

func TestPrintWorkingDirectory(t *testing.T) {
	shell := newTestShell()

	var out bytes.Buffer
	shell.cmd.SetOut(&out)

	shell.printWorkingDirectory()

	currDir, err := os.Getwd()
	assert.NoError(t, err)

	assert.Contains(t, out.String(), currDir)
}

func TestListFiles(t *testing.T) {
	shell := newTestShell()

	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "dir1"), 0755)

	shell.changeDirectory([]string{tmpDir})

	files, _ := shell.listFiles()

	assert.Contains(t, files, "file1.txt")
	assert.Contains(t, files, "dir1/")
}
