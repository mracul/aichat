package interfaces
package interfaces

type Context interface {
	App() interface{}
	GUI() interface{}
	Storage() interface{}
	Config() interface{}
	Logger() interface{}
}

