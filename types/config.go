package types

// config.go - Application configuration types and helpers for the unified app model
// Contains AppConfig and related configuration logic for the app package.

import (
	"gopkg.in/ini.v1"
)

const settingsPath = ".config/settings.ini"

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

// GetCurrentTheme reads the current theme from the settings.ini file
func GetCurrentTheme() (string, error) {
	cfg, err := ini.Load(settingsPath)
	if err != nil {
		return "Default", nil // fallback to Default if not found
	}
	return cfg.Section("Theme").Key("currentTheme").String(), nil
}

// SetCurrentTheme writes the current theme to the settings.ini file
func SetCurrentTheme(themeName string) error {
	cfg, err := ini.LoadSources(ini.LoadOptions{Loose: true}, settingsPath)
	if err != nil {
		cfg = ini.Empty()
	}
	cfg.Section("Theme").Key("currentTheme").SetValue(themeName)
	return cfg.SaveTo(settingsPath)
}
