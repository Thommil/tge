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
	app            App
	ticker         *time.Ticker
	canvas         js.Value
	isPaused       bool
	isStopped      bool
	isPausedChan   chan bool
	isStoppedChan  chan bool
	isDisposedChan chan bool
}

func (runtime browserRuntime) Stop() {
	go func() {
		runtime.ticker.Stop()
		runtime.isPausedChan <- true
		runtime.isStoppedChan <- true
	}()
}

func doRun(app App, settings *Settings) error {
	log.Println("doRun()")

	// -------------------------------------------------------------------- //
	// Init
	// -------------------------------------------------------------------- //
	jsTge := js.Global().Get("tge")
	jsTge.Call("resize", settings.Width, settings.Height)
	jsTge.Call("setFullscreen", settings.Fullscreen)

	canvas := jsTge.Call("init")

	// Instanciate Runtime
	browserRuntime := browserRuntime{
		app:            app,
		isPaused:       true,
		isStopped:      false,
		canvas:         canvas,
		isPausedChan:   make(chan bool),
		isStoppedChan:  make(chan bool),
		isDisposedChan: make(chan bool),
	}

	// Start App
	app.OnStart(&browserRuntime)
	app.OnResume()
	browserRuntime.isPaused = false
	app.OnResize(browserRuntime.canvas.Get("clientWidth").Int(),
		browserRuntime.canvas.Get("clientHeight").Int())

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	tpsDelay := time.Duration(1000000000 / settings.TPS)
	browserRuntime.ticker = time.NewTicker(tpsDelay)
	defer browserRuntime.ticker.Stop()

	mutex := &sync.Mutex{}
	elapsedTpsTime := time.Duration(0)
	go func() {
		for range browserRuntime.ticker.C {
			if !browserRuntime.isPaused {
				startTps := time.Now()
				app.OnTick(elapsedTpsTime, mutex)
				elapsedTpsTime = (tpsDelay - time.Since(startTps))
				time.Sleep(elapsedTpsTime)
			}
		}
	}()

	// -------------------------------------------------------------------- //
	// Callbacks
	// -------------------------------------------------------------------- //

	// window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
	// 	// OS Specific - Windows call resize to 0
	// 	if !desktopRuntime.isPaused && width > 0 {
	// 		app.OnResize(width, height)
	// 	}
	// })

	// window.SetFocusCallback(func(w *glfw.Window, focused bool) {
	// 	if focused && desktopRuntime.isPaused {
	// 		desktopRuntime.isPaused = false
	// 		app.OnResume()
	// 		// OS Specific - MacOS do not resize at start
	// 		resizeAtStart.Do(func() {
	// 			if runtime.GOOS != "windows" {
	// 				app.OnResize(settings.Width, settings.Height)
	// 			}
	// 		})
	// 	} else if !desktopRuntime.isPaused {
	// 		desktopRuntime.isPaused = true
	// 		app.OnPause()
	// 	}
	// })

	// window.SetCloseCallback(func(w *glfw.Window) {
	// 	desktopRuntime.ticker.Stop()
	// 	if !desktopRuntime.isPaused {
	// 		app.OnPause()
	// 	}
	// 	app.OnStop()
	// })

	// -------------------------------------------------------------------- //
	// Render Loop
	// -------------------------------------------------------------------- //
	fpsDelay := time.Duration(1000000000 / settings.FPS)

	var renderFrame js.Callback
	defer renderFrame.Release()
	elapsedFpsTime := time.Duration(0)
	renderFrame = js.NewCallback(func(args []js.Value) {
		// Get channels to check status
		select {
		case browserRuntime.isPaused = <-browserRuntime.isPausedChan:
			if browserRuntime.isPaused {
				browserRuntime.app.OnPause()
			}
		case browserRuntime.isStopped = <-browserRuntime.isStoppedChan:
			if browserRuntime.isStopped {
				browserRuntime.app.OnStop()
				browserRuntime.isDisposedChan <- true
			}
		default:
		}

		// Render
		if !browserRuntime.isStopped {
			if !browserRuntime.isPaused {
				startFps := time.Now()
				app.OnRender(elapsedFpsTime, mutex)
				elapsedFpsTime = (fpsDelay - time.Since(startFps))
				time.Sleep(elapsedFpsTime)
				js.Global().Call("requestAnimationFrame", renderFrame)
			} else {
				time.Sleep(fpsDelay)
				js.Global().Call("requestAnimationFrame", renderFrame)
			}
		}
	})
	js.Global().Call("requestAnimationFrame", renderFrame)

	// Block until dispose chan notified
	<-browserRuntime.isDisposedChan

	// Call JS stop handler
	jsTge.Call("stop")

	return nil
}
