package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/balena-os/balena-extension-runtime/internal/log"
	"github.com/balena-os/balena-extension-runtime/internal/manager"
	"github.com/balena-os/balena-extension-runtime/internal/version"
	"github.com/spf13/cobra"
)

var (
	logLevel string
	logger   *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:     "balena-extension-manager",
	Short:   "Manage hostapp extension lifecycle",
	Version: fmt.Sprintf("%s (commit: %s)", version.Version, version.GitCommit),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		level, err := parseLogLevel(logLevel)
		if err != nil {
			return err
		}
		logger = log.NewLogger(level)
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Re-create extension containers for the new kernel version",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootfs, _ := cmd.Flags().GetString("rootfs")
		return manager.Update(context.Background(), logger, rootfs)
	},
}

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Remove stale extension containers and orphaned images",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manager.Cleanup(context.Background(), logger)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info",
		"Set the logging level (debug, info, warn, error)")
	updateCmd.Flags().String("rootfs", "", "Path to the new OS rootfs")
	updateCmd.MarkFlagRequired("rootfs")
	rootCmd.AddCommand(updateCmd, cleanupCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func parseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("invalid log level %q", level)
	}
}
