// +build js

package tge

import (
	fmt "fmt"
	sync "sync"
	js "syscall/js"
	time "time"
)

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type browserRuntime struct {
	app          App
	plugins      []Plugin
	ticker       *time.Ticker
	canvas       js.Value
	isPaused     bool
	isStopped    bool
	tickTicker   *time.Ticker
	renderTicker *time.Ticker
}

func (runtime *browserRuntime) Use(plugin Plugin) {
	runtime.plugins = append(runtime.plugins, plugin)
	err := plugin.Init(runtime)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (runtime *browserRuntime) Stop() {
	if !runtime.isPaused {
		runtime.isPaused = true
		runtime.app.OnPause()
	}
	runtime.isStopped = true
	runtime.app.OnStop()
	runtime.app.OnDispose()
}

// Run main entry point of runtime
func Run(app App) error {
	fmt.Println("Run()")

	// -------------------------------------------------------------------- //
	// Create
	// -------------------------------------------------------------------- //
	settings := &defaultSettings
	err := app.OnCreate(settings)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// -------------------------------------------------------------------- //
	// Init
	// -------------------------------------------------------------------- //
	jsTge := js.Global().Get("tge")
	if settings.Fullscreen {
		jsTge.Call("setFullscreen", settings.Fullscreen)
	} else {
		jsTge.Call("resize", settings.Width, settings.Height)
	}

	canvas := jsTge.Call("init")

	// Instanciate Runtime
	browserRuntime := &browserRuntime{
		app:       app,
		plugins:   make([]Plugin, 0),
		isPaused:  true,
		isStopped: true,
		canvas:    canvas,
	}

	// Start App
	app.OnStart(browserRuntime)
	browserRuntime.isStopped = false

	// Resume App
	app.OnResume()
	browserRuntime.isPaused = false

	// Resize App
	app.OnResize(browserRuntime.canvas.Get("clientWidth").Int(),
		browserRuntime.canvas.Get("clientHeight").Int())

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	tpsDelay := time.Duration(1000000000 / settings.TPS)
	browserRuntime.tickTicker = time.NewTicker(tpsDelay)
	defer browserRuntime.tickTicker.Stop() // Avoid leak

	mutex := &sync.Mutex{}
	elapsedTpsTime := time.Duration(0)
	go func() {
		for now := range browserRuntime.tickTicker.C {
			if !browserRuntime.isPaused {
				app.OnTick(elapsedTpsTime, mutex)
				elapsedTpsTime = tpsDelay - time.Since(now)
				if elapsedTpsTime < 0 {
					elapsedTpsTime = 0
				}
			} else if browserRuntime.isStopped {
				browserRuntime.tickTicker.Stop()
				return
			}
		}
	}()

	// -------------------------------------------------------------------- //
	// Callbacks
	// -------------------------------------------------------------------- //

	// Resize
	js.Global().Call("addEventListener", "resize", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped {
			app.OnResize(browserRuntime.canvas.Get("clientWidth").Int(),
				browserRuntime.canvas.Get("clientHeight").Int())
		}
	}))

	// Focus
	browserRuntime.canvas.Call("addEventListener", "blur", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped && !browserRuntime.isPaused {
			go func() {
				browserRuntime.isPaused = true
				browserRuntime.app.OnPause()
			}()
		}
	}))

	browserRuntime.canvas.Call("addEventListener", "focus", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped && browserRuntime.isPaused {
			go func() {
				browserRuntime.app.OnResume()
				browserRuntime.isPaused = false
			}()
		}
	}))

	// Destroy
	js.Global().Call("addEventListener", "beforeunload", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped {
			browserRuntime.Stop()
		}
	}))

	// -------------------------------------------------------------------- //
	// Render Loop
	// -------------------------------------------------------------------- //
	// -------------------------------------------------------------------- //
	// Render Loop
	// -------------------------------------------------------------------- //
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	elapsedFpsTime := time.Duration(0)
	browserRuntime.renderTicker = time.NewTicker(fpsDelay)
	defer browserRuntime.renderTicker.Stop() // Avoid leak

	renderFrame := js.NewCallback(func(args []js.Value) {
		now := time.Now()
		app.OnRender(elapsedFpsTime, mutex)
		elapsedFpsTime = fpsDelay - time.Since(now)
		if elapsedFpsTime < 0 {
			elapsedFpsTime = 0
		}
	})

	for range browserRuntime.renderTicker.C {
		if !browserRuntime.isPaused {
			js.Global().Call("requestAnimationFrame", renderFrame)
		} else if browserRuntime.isStopped {
			browserRuntime.renderTicker.Stop()
			jsTge.Call("stop")
			return nil
		}
	}

	return nil
}
