package tge

import (
	runtime "github.com/thommil/tge/runtime"
)

// App defines API to implement for TGE applications
type App interface {
	Create(settings Settings) error
	Start(runtime runtime.Runtime) error
	Resize(width int, height int) error
	Resume() error
	Pause() error
	Dispose() error
}

// Instanciate is the main entry point
func Instanciate(app App) {
	app.Create(Settings{})
	//app.Start()
}
