package main

import (
	"time"

	// Runtime package
	tge "github.com/thommil/tge"
	//Available plugins - Just uncomment to enable
	//gesture "github.com/thommil/tge-gesture"
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
	// Settings is passed as a pointer can be modified to adapt App to your needs
	// Ex :
	//	settings.Name = "My Awesome App++"
	//	settings.Fullscreen = true
	//	settings.EventMask = tge.AllEventsEnable
	// 	...
	return nil
}

// OnStart is called when all Runtime resources are available but looping as not been
// started. Initializations should be done here (GL conf, physics engine ...).
func (app *App) OnStart(runtime tge.Runtime) error {
	// Here is the place where your initialize everything, the runtime is UP but the loops
	// are not started :
	// 	- OpenGL
	//	- Physics
	//	- AI
	//	- Load previous Save
	//	- State Machine setup
	// 	...
	return nil
}

// OnResume is called just after OnStart and also when the Runtime is awaken after a pause.
func (app *App) OnResume() {
	// This method is called to notify that loops are started, don't place heavey treatments here
	// Most of the time, this method is just used to hide/remove a pause screen
}

// OnRender is called each time the graphical context is redrawn. This method should only implements
// graphical calls and not logical ones. The syncChan is used to wait for Tick() loop draw commands
func (app *App) OnRender(elaspedTime time.Duration, syncChan <-chan interface{}) {
	// Always listen at least for one object from syncChan to synchronize it with Tick(), in other case the Tick() loop
	// will be blocked
	<-syncChan

	// This loop is dedicated to graphical/GPU treatments:
	//	- OpenGl Calls
	//	- Vulkan Calls (one day ;)
	//
	// Data becomes available from syncChan pipe and sent by the Tick loop, as mentioned below, it's possible to
	// reused syncChan several times in a single Render/Tick call for progressive rendering.
	//
	// As syncChan is a generic interface channel, it's also possible to select treatment to apply:
	//	data := <-syncChan
	//	switch data.(type) {
	//		...
	//	}
}

// OnTick handles logical operations like physics, AI or any background task not relatd to graphics.
// The syncChan is used to notify Render() loop with draw commands.
func (app *App) OnTick(elaspedTime time.Duration, syncChan chan<- interface{}) {
	// This loop is dedicated to logical/CPU treatments:
	//	- physics
	//	- AI
	//	- State Machine
	//	- File access ...
	//
	// Each time Tick loop needs to send data to Render loop, use the syncChan.
	//
	// In can be done once per call or several times if you want a progressive rendering
	//
	// As data can be shared between Tick and Render loops, a good practice is too handle heavy treatments
	// in Tick dedicated data, then copy data to Render dedicated data and send it through the syncChan.
	//
	// A good candidate for copy if the reflect.Copy() function:
	//   reflect.Copy(reflect.ValueOf(renderData), reflect.ValueOf(tickData))
	//
	// If your data is based on something else than slices but its size justifies low level memory copy, you can
	// also put ticker data in single element slice and use reflect.Copy().
	//
	// Tick loop is running in a dedicated Go routine, it's also possible to start subroutines in this loop to
	// benefit from available cores and increase treatment speed using map/reduce oriented algorithms.
	//
	// Always send something to syncChan to synchronize it with Render(), in other case the Tick() loop
	// will be a simple infinite loop and your App will destroy your Desktop/Mobile/Browser
	syncChan <- true
}

// OnPause is called when the Runtime lose focus (alt-tab, home button, tab change ...) This is a good entry point to
// set and display a pause screen
func (app *App) OnPause() {
	// Most of the time, just set a flag indicating the paused state of your App
}

// OnStop is called when the Runtime is ending, context saving should be done here. On current Android version this
// handler is also called when the application is paused (and restart after).
func (app *App) OnStop() {
	// This is where you backup everything if needed (state machine, save ...) The runtime tries to call and execute
	// this method before leaving to allow proper exit but nothing is guaranteed on some targets (WEB)
}

// OnDispose is called when all exit treatments are done for cleaning task (memory, tmp files ...)
func (app *App) OnDispose() {
	// Optional but always good practice to clean up everything before leaving :)
}

// Main entry point, simply instanciates App and runs it through Runtime
func main() {
	// The line below should be the only one, code here is not portable!
	tge.Run(&App{})
}
