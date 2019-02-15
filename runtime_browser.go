// +build js

package tge

import (
	log "log"
	sync "sync"
	js "syscall/js"
	time "time"
)

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type browserRuntime struct {
	app          App
	ticker       *time.Ticker
	canvas       js.Value
	isPaused     bool
	isStopped    bool
	isPausedChan chan bool
	tickTicker   *time.Ticker
	tickEnd      chan bool
	renderTicker *time.Ticker
	renderEnd    chan bool
}

func (runtime browserRuntime) Stop() {
	runtime.isPaused = true
	go func() {
		runtime.isPausedChan <- true
		runtime.app.OnPause()
		go func() {

			runtime.tickEnd <- true
			runtime.renderEnd <- true
		}()
	}()

}

// Run main entry point of runtime
func Run(app App) error {
	log.Println("Run()")

	// -------------------------------------------------------------------- //
	// Create
	// -------------------------------------------------------------------- //
	settings := &defaultSettings
	err := app.OnCreate(settings)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.OnDispose() // Should be never called

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
		app:          app,
		isPaused:     true,
		isStopped:    false,
		canvas:       canvas,
		isPausedChan: make(chan bool),
		tickEnd:      make(chan bool),
		renderEnd:    make(chan bool),
	}

	// Start App
	app.OnStart(browserRuntime)
	app.OnResume()
	browserRuntime.isPaused = false
	app.OnResize(browserRuntime.canvas.Get("clientWidth").Int(),
		browserRuntime.canvas.Get("clientHeight").Int())

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	tpsDelay := time.Duration(1000000000 / settings.TPS)
	browserRuntime.tickTicker = time.NewTicker(tpsDelay)
	defer browserRuntime.tickTicker.Stop()

	mutex := &sync.Mutex{}
	elapsedTpsTime := time.Duration(0)
	go func() {
		for {
			select {
			case <-browserRuntime.tickEnd:
				return
			case now := <-browserRuntime.tickTicker.C:
				if !browserRuntime.isPaused {
					app.OnTick(elapsedTpsTime, mutex)
					elapsedTpsTime = tpsDelay - time.Since(now)
					if elapsedTpsTime < 0 {
						elapsedTpsTime = 0
					}
				}
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
				browserRuntime.isPausedChan <- true
				browserRuntime.app.OnPause()
			}()
		}
	}))

	browserRuntime.canvas.Call("addEventListener", "focus", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped && browserRuntime.isPaused {
			go func() {
				browserRuntime.app.OnResume()
				browserRuntime.isPausedChan <- false
			}()
		}
	}))

	// Destroy
	js.Global().Call("addEventListener", "beforeunload", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped {
			app.OnStop()
			app.OnDispose()
		}
	}))

	// -------------------------------------------------------------------- //
	// Render Loop
	// -------------------------------------------------------------------- //
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	now := time.Now()
	elapsedFpsTime := time.Duration(0)
	browserRuntime.renderTicker = time.NewTicker(fpsDelay)
	defer browserRuntime.renderTicker.Stop()

	renderFrame := js.NewCallback(func(args []js.Value) {
		now = time.Now()
		app.OnRender(elapsedFpsTime, mutex)
		elapsedFpsTime = fpsDelay - time.Since(now)
		if elapsedFpsTime < 0 {
			elapsedFpsTime = 0
		}
	})

	for {
		select {
		case <-browserRuntime.renderEnd:
			browserRuntime.tickTicker.Stop()
			browserRuntime.renderTicker.Stop()
			renderFrame.Release()
			browserRuntime.isStopped = true
			app.OnStop()
			jsTge.Call("stop")
			app.OnDispose()
			<-make(chan int)
		case <-browserRuntime.renderTicker.C:
			select {
			case browserRuntime.isPaused = <-browserRuntime.isPausedChan:
			default:
			}
			if !browserRuntime.isPaused {
				js.Global().Call("requestAnimationFrame", renderFrame)
			}
		}
	}

}
