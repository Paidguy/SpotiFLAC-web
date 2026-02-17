package backend

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
)

func OpenFolderInExplorer(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	return cmd.Start()
}

// SelectFolderDialog is not supported in web mode
func SelectFolderDialog(ctx context.Context, defaultPath string) (string, error) {
	return "", fmt.Errorf("folder selection dialogs are not supported in web mode")
}

// SelectFileDialog is not supported in web mode
func SelectFileDialog(ctx context.Context) (string, error) {
	return "", fmt.Errorf("file selection dialogs are not supported in web mode")
}

// SelectImageVideoDialog is not supported in web mode
func SelectImageVideoDialog(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("file selection dialogs are not supported in web mode")
}
