// +build js

package tge

import (
	log "log"
	"time"

	physics "github.com/thommil/tge/physics"
	player "github.com/thommil/tge/player"
	renderer "github.com/thommil/tge/renderer"
	ui "github.com/thommil/tge/ui"
)

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type browserRuntime struct {
	app App

	renderer renderer.Renderer
	ui       ui.UI
	player   player.Player
	physics  physics.Physics
	ticker   *time.Ticker
	isPaused bool
}

func (runtime browserRuntime) GetRenderer() renderer.Renderer {
	return runtime.renderer
}

func (runtime browserRuntime) GetUI() ui.UI {
	return runtime.ui
}

func (runtime browserRuntime) GetPlayer() player.Player {
	return runtime.player
}

func (runtime browserRuntime) GetPhysics() physics.Physics {
	return runtime.physics
}

func (runtime browserRuntime) Stop() {
	runtime.ticker.Stop()
	runtime.app.OnPause()
	runtime.app.OnStop()
}

func doRun(app App, settings *Settings) error {
	log.Println("doRun()")

	// -------------------------------------------------------------------- //
	// Init
	// -------------------------------------------------------------------- //

	// Instanciate Runtime
	browserRuntime := browserRuntime{
		app:      app,
		isPaused: true,
	}

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	tpsDelay := time.Duration(1000000000 / settings.Physics.TPS)
	browserRuntime.ticker = time.NewTicker(tpsDelay)
	defer browserRuntime.ticker.Stop()
	go func() {
		for range browserRuntime.ticker.C {
			if !browserRuntime.isPaused {
				app.OnTick(tpsDelay)
			}
		}
	}()

	return nil
}
