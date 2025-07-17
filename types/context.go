package types

// AppProvider provides App
type AppProvider interface {
	App() interface{}
}

// GUIProvider provides GUI
type GUIProvider interface {
	GUI() interface{}
}

// StorageProvider provides Storage
type StorageProvider interface {
	Storage() interface{}
}

// ConfigProvider provides Config
type ConfigProvider interface {
	Config() interface{}
}

// LoggerProvider provides Logger
type LoggerProvider interface {
	Logger() interface{}
}

// Context composes all providers
type Context interface {
	AppProvider
	GUIProvider
	StorageProvider
	ConfigProvider
	LoggerProvider
}
