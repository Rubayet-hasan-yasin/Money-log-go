package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rubayet-hasan-yasin/sceleton-go/internal/config"
	"github.com/rubayet-hasan-yasin/sceleton-go/internal/infrastructure/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serveCmd = &cobra.Command{
	Use:   "serve-rest",
	Short: "Start the HTTP REST server",
	Long:  `Start the REST API server with all endpoints, database connection, and middleware.`,
	Example: `  # Start server with default settings
  go run ./cmd serve-rest
  
  # Start server with custom port
  go run ./cmd serve-rest --port 9090`,
	RunE: runServe,
}

var (
	port string
)

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", "", "server port (overrides config)")
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg := config.GetConfig()

	if port != "" {
		cfg.Server.Port = port
	}

	// Initialize logger
	log, err := logger.NewLogger(cfg.Server.AppEnv)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer log.Sync()

	// Create context with signal cancellation
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run HTTP server
	if err := runHTTPServer(ctx, cfg, log); err != nil {
		log.Error("server exited with error", zap.Error(err))
		return err
	}

	log.Info("Server shutdown complete")
	return nil
}
