// main.go - Entry point for the AI CLI application
// This file handles environment setup, configuration loading, and launches the unified application
// The unified app model integrates responsive resizing, ANSI optimization, and proper navigation

package main

import (
	"aichat/app"
	"aichat/services/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

// =====================================================================================
// 🚀 Application Entry Point
// =====================================================================================

func main() {
	// Initialize logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Starting AI CLI application", "version", "1.0.0")

	// Create navigation storage directly
	navStorage := storage.NewNavigationStorage(".config")

	// Create and run the unified application (optimized mode by default)
	cfg := config.DefaultAppConfig()
	_ = cfg // Currently not used by NewAppModel, but available for future use
	appModel := app.NewAppModel(navStorage, nil, "")
	program := tea.NewProgram(appModel, tea.WithAltScreen(), tea.WithMouseCellMotion())

	// Set up graceful shutdown
	setupGracefulShutdown(program, logger)

	// Run the application
	if _, err := program.Run(); err != nil {
		logger.Error("Application failed", "error", err)
		os.Exit(1)
	}

	logger.Info("Application completed successfully")
}

// =====================================================================================
// 🛡️ Graceful Shutdown
// =====================================================================================

// setupGracefulShutdown sets up signal handling for graceful shutdown
func setupGracefulShutdown(program *tea.Program, logger *slog.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Received shutdown signal, cleaning up...")

		// Note: tea.Program doesn't expose Model() method directly
		// The shutdown will be handled by the program's cleanup
		program.Quit()
	}()
}
