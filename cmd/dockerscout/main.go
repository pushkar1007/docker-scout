// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (c) 2024-2026 dockerscout contributors
// https://github.com/fr4nsys/dockerscout

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/fr4nsys/dockerscout/internal/app"
)

var (
	cfgFile string
	mode    string
)

var rootCmd = &cobra.Command{
	Use:   "dockerscout",
	Short: "Docker Management Platform",
	Long:  `dockerscout is a self-hosted Docker management platform with security scoring, centralized config, and NPM integration.`,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  `Start the dockerscout server in the specified mode (standalone, master, or agent).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.Run(cfgFile, mode)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunMigrations(cfgFile, "up")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down [N]",
	Short: "Rollback N migrations (default: 1)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		steps := "1"
		if len(args) > 0 {
			steps = args[0]
		}
		return app.RunMigrations(cfgFile, "down:"+steps)
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return app.RunMigrations(cfgFile, "status")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		app.PrintVersion()
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
}

var configCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := app.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("configuration error: %w", err)
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("validation error: %w", err)
		}
		fmt.Println("Configuration is valid")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration (sensitive values masked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := app.LoadConfig(cfgFile)
		if err != nil {
			return err
		}
		cfg.PrintMasked()
		return nil
	},
}

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin management commands",
}

var adminResetPasswordCmd = &cobra.Command{
	Use:   "reset-password [NEW_PASSWORD]",
	Short: "Reset admin user password (or create admin if missing)",
	Long: `Reset the password for the 'admin' user. If the admin user
doesn't exist, it will be created. Also unlocks the account
if it was locked due to failed login attempts.

If no password is provided, defaults to 'dockerscout'.

Usage from docker:
  docker exec dockerscout-app /app/dockerscout admin reset-password MyNewPass123`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		password := "dockerscout"
		if len(args) > 0 {
			password = args[0]
		}
		if len(password) < 8 {
			return fmt.Errorf("password must be at least 8 characters")
		}
		return app.ResetAdminPassword(cfgFile, password)
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default: /etc/dockerscout/config.yaml or ./config.yaml)")

	// Serve flags
	serveCmd.Flags().StringVarP(&mode, "mode", "m", "standalone", "operation mode: standalone|master|agent")
	serveCmd.Flags().String("component", "", "component to run: api|gateway|scheduler (master mode only)")

	// Build command tree
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	rootCmd.AddCommand(migrateCmd)

	configCmd.AddCommand(configCheckCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)

	adminCmd.AddCommand(adminResetPasswordCmd)
	rootCmd.AddCommand(adminCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
