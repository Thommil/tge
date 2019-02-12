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
	Create(settings *Settings) error
	Start() error
	Resize(width int, height int)
	Resume()
	Render(renderer renderer.Renderer, ui ui.UI, player player.Player)
	Tick(physics physics.Physics)
	Pause()
	Stop()
	Dispose() error
}

// Runtime API
type Runtime interface {
}

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

// Run is the main entry point
func Run(app App) {
	log.Println("Run()")

	settings := Settings{}
	app.Create(&settings)

	err := doRun(app, &settings)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.Dispose()
}
