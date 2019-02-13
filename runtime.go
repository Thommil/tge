package tge

import (
	log "log"
	"time"

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
	OnRender(elaspedTime time.Duration)
	OnTick(elaspedTime time.Duration)
	OnPause()
	OnStop()
	OnDispose() error
}

// Runtime API
type Runtime interface {
	GetRenderer() renderer.Renderer
	GetUI() ui.UI
	GetPlayer() player.Player
	GetPhysics() physics.Physics
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
