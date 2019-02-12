package tge

import (
	log "log"

	physics "github.com/thommil/tge/physics"
	player "github.com/thommil/tge/player"
	renderer "github.com/thommil/tge/renderer"
	ui "github.com/thommil/tge/ui"
)

// App defines API to implement for TGE applications
type App interface {
	OnCreate(settings *Settings) error
	OnStart(runtime Runtime) error
	OnResize(width int, height int)
	OnResume()
	OnRender(renderer renderer.Renderer, ui ui.UI, player player.Player)
	OnTick(physics physics.Physics)
	OnPause()
	OnStop()
	OnDispose() error
}

// Runtime API
type Runtime interface {
	Stop()
}

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

// Run is the main entry point
func Run(app App) {
	log.Println("Run()")

	settings := &defaultSettings
	app.OnCreate(settings)

	err := doRun(app, settings)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.OnDispose()
}
