package hooks

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/balena-os/balena-extension-runtime/internal/labels"
)

// ExecuteIfPresent runs a hook script from the extension rootfs if it exists.
// The hook path is relative to rootfs (e.g., "hooks/create").
// Returns nil if the hook does not exist.
func ExecuteIfPresent(logger *slog.Logger, rootfs string, hookPath string, annotations map[string]string) error {
	absPath := filepath.Join(rootfs, hookPath)

	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		logger.Debug("hook not present, skipping", "hook", hookPath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to stat hook %s: %w", absPath, err)
	}
	if info.IsDir() {
		return fmt.Errorf("hook %s is a directory, not executable", absPath)
	}
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("hook %s is not executable", absPath)
	}

	logger.Info("executing hook", "hook", hookPath, "rootfs", rootfs)

	env := append(os.Environ(),
		"EXTENSION_ROOTFS="+rootfs,
	)
	env = append(env, labels.ToEnv(annotations)...)

	cmd := exec.Command(absPath)
	cmd.Env = env
	cmd.Stdout = os.Stderr // hooks log to runtime stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("hook %s failed: %w", hookPath, err)
	}

	logger.Info("hook completed", "hook", hookPath)
	return nil
}
