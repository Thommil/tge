package main

import (
	"sync"
	"time"

	// Runtime package
	tge "github.com/thommil/tge"
	//Available plugins - Just uncomment to enable
	//gl "github.com/thommil/tge-gl"
	//g3n "github.com/thommil/tge-g3n"
)

// App instance context and definition
type App struct {
	// Put global attributes of your application here
}

// OnCreate is called at App instanciation, the Runtime resources are not
// yet available. Settings and treatments not related to Runtime should be done here.
func (app *App) OnCreate(settings *tge.Settings) error {

	return nil
}

// OnStart is called when all Runtime resources are available but looping as not been
// started. Initializations should be done here (GL conf, physics engine ...).
func (app *App) OnStart(runtime tge.Runtime) error {

	return nil
}

// OnResume is called just after OnStart and also when the Runtime is awaken after a pause.
func (app *App) OnResume() {

}

// OnRender is called each time the graphical context is redrawn. This method should only implements
// graphical calls and not logical ones. The mutex allows to synchronize critical path with Tick() loop.
func (app *App) OnRender(elaspedTime time.Duration, mutex *sync.Mutex) {

}

// OnTick is called at a rate defined in settings, logical operations should be done here like physics, AI
// or any background task not relatd to graphics. The mutex allows to synchronize critical path with Render() loop.
func (app *App) OnTick(elaspedTime time.Duration, mutex *sync.Mutex) {

}

// OnPause is called when the Runtime lose focus (alt-tab, home button, tab change ...) This is a good nrty point to
// set and display a pause screen
func (app *App) OnPause() {

}

// OnStop is called when the Runtime is ending, context saving should be done here. On current Android version this
// handler is also called when the application is paused (and restart after).
func (app *App) OnStop() {

}

// OnDispose is called when all exit treatments are done for cleaning task (memory, tmp files ...)
func (app *App) OnDispose() {

	return nil
}

// Main entry point, simply instanciates App and runs it through Runtime
func main() {
	tge.Run(&App{})
}
