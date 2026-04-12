package runtime

import (
	"fmt"
	"syscall"

	"github.com/balena-os/balena-extension-runtime/internal/oci"
	"github.com/balena-os/balena-extension-runtime/internal/proxy"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// Kill sends a signal to the proxy process.
func Kill(containerID string, signal syscall.Signal) error {
	state, err := oci.ReadState(containerID)
	if err != nil {
		return fmt.Errorf("failed to read state: %w", err)
	}

	if state.Pid > 0 {
		if err := proxy.Signal(state.Pid, signal); err != nil {
			// Process may already be dead — not an error for extensions
			if signal == syscall.SIGKILL || signal == syscall.SIGTERM {
				// Best effort
			} else {
				return fmt.Errorf("failed to send signal: %w", err)
			}
		}
	}

	state.Status = specs.StateStopped
	if err := oci.WriteState(state); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}
	return nil
}
