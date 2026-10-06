package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Go Backend Starter API Server",
	Long: `A modular, scalable, clean architecture Go backend starter skeleton.
Designed as a robust foundation for building modern APIs with authentication,
worker queues, PostgreSQL persistence, S3 storage, and REST APIs.`,
	Version: "1.0.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", ".env", "config file path")
}
