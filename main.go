package main
package main

import (
	"aichat/types"
	"aichat/services/storage"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	logger.Info("Starting AI CLI application", "version", "1.0.0")

	navStorage := storage.NewNavigationStorage(".config")
	cfg := app.DefaultAppConfig()
	appModel := app.NewUnifiedAppModel(cfg, navStorage, logger)

	program := tea.NewProgram(appModel, tea.WithAltScreen(), tea.WithMouseCellMotion())
	setupGracefulShutdown(program, logger)

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

