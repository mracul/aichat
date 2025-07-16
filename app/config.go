// config.go - Application configuration types and helpers for the unified app model
// Contains AppConfig and related configuration logic for the app package.

package app

type AppConfig struct {
	EnableCaching bool
	EnableLogging bool
	DefaultWidth  int
	DefaultHeight int
	MinWidth      int
	MinHeight     int
}

// DefaultAppConfig returns default application configuration
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		EnableCaching: true,
		EnableLogging: true,
		DefaultWidth:  80,
		DefaultHeight: 24,
		MinWidth:      40,
		MinHeight:     10,
	}
}
