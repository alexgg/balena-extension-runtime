package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var proxyContainerID string

var proxyCmd = &cobra.Command{
	Use:    "proxy",
	Short:  "Proxy process that provides a PID for the containerd shim",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("proxy started", "container", proxyContainerID)

		// Block until SIGUSR1 (start complete) or SIGTERM (killed).
		// For extensions, SIGUSR1 means "start finished" — exit cleanly.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sigCh
		switch sig {
		case syscall.SIGUSR1:
			// Extension started — exit cleanly so container becomes "Exited (0)"
			os.Exit(0)
		case syscall.SIGTERM, syscall.SIGINT:
			os.Exit(0)
		}

		return nil
	},
}

func init() {
	proxyCmd.Flags().StringVar(&proxyContainerID, "id", "", "Container ID")
	rootCmd.AddCommand(proxyCmd)
}
