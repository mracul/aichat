package navigation
package navigation

import (
	"aichat/services/storage"
	"log/slog"
)

type AppContext struct {
	app     interface{} // *UnifiedAppModel
	gui     interface{} // *GUIAppModel, nil for TUI
	storage storage.NavigationStorage
	config  interface{} // *AppConfig
	logger  *slog.Logger
}

func NewContext(app interface{}, gui interface{}, storage storage.NavigationStorage, config interface{}, logger *slog.Logger) *AppContext {
	return &AppContext{
		app:     app,
		gui:     gui,
		storage: storage,
		config:  config,
		logger:  logger,
	}
}

func (ctx *AppContext) App() interface{}                   { return ctx.app }
func (ctx *AppContext) GUI() interface{}                   { return ctx.gui }
func (ctx *AppContext) Storage() storage.NavigationStorage { return ctx.storage }
func (ctx *AppContext) Config() interface{}                { return ctx.config }
func (ctx *AppContext) Logger() *slog.Logger               { return ctx.logger }
