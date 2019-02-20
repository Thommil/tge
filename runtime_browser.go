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
	app       App
	plugins   map[string]Plugin
	ticker    *time.Ticker
	canvas    *js.Value
	isPaused  bool
	isStopped bool
	done      chan bool
}

func (runtime *browserRuntime) Use(plugin Plugin) {
	name := plugin.GetName()
	fmt.Printf("Loading plugin %s\n", name)
	runtime.plugins[name] = plugin
	err := plugin.Init(runtime)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (runtime *browserRuntime) GetPlugin(name string) Plugin {
	return runtime.plugins[name]
}

func (runtime *browserRuntime) GetHost() interface{} {
	host := js.Global()
	return &host
}

func (runtime *browserRuntime) GetRenderer() interface{} {
	glContext := runtime.canvas.Call("getContext", "webgl")
	if glContext == js.Undefined() {
		glContext = runtime.canvas.Call("getContext", "experimental-webgl")
	}
	if glContext == js.Undefined() {
		err := fmt.Errorf("No WebGL support found in brower")
		fmt.Println(err)
		panic(err)
	}
	return &glContext
}

func (runtime *browserRuntime) Stop() {
	if !runtime.isPaused {
		runtime.isPaused = true
		runtime.app.OnPause()
	}
	runtime.isStopped = true
	runtime.app.OnStop()
	// Unload plugins
	for _, plugin := range runtime.plugins {
		plugin.Dispose()
	}
	runtime.app.OnDispose()
}

// Run main entry point of runtime
func Run(app App) error {
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
		plugins:   make(map[string]Plugin),
		isPaused:  true,
		isStopped: true,
		canvas:    &canvas,
		done:      make(chan bool),
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
	mutex := &sync.Mutex{}
	tpsDelay := time.Duration(1000000000 / settings.TPS)
	elapsedTpsTime := time.Duration(0)
	go func() {
		for !browserRuntime.isStopped {
			if !browserRuntime.isPaused {
				now := time.Now()
				app.OnTick(elapsedTpsTime, mutex)
				elapsedTpsTime = tpsDelay - time.Since(now)
				if elapsedTpsTime < 0 {
					elapsedTpsTime = 0
				}
				time.Sleep(elapsedTpsTime)
			} else {
				time.Sleep(tpsDelay)
			}
		}
	}()

	// -------------------------------------------------------------------- //
	// Callbacks
	// -------------------------------------------------------------------- //

	// Resize
	js.Global().Call("addEventListener", "resize", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped && !browserRuntime.isPaused {
			app.OnResize(browserRuntime.canvas.Get("clientWidth").Int(),
				browserRuntime.canvas.Get("clientHeight").Int())
		}
	}))

	// Focus
	browserRuntime.canvas.Call("addEventListener", "blur", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped && !browserRuntime.isPaused {
			browserRuntime.isPaused = true
			browserRuntime.app.OnPause()
		}
	}))

	browserRuntime.canvas.Call("addEventListener", "focus", js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isStopped && browserRuntime.isPaused {
			browserRuntime.app.OnResume()
			browserRuntime.isPaused = false
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
	var renderFrame js.Callback
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	elapsedFpsTime := time.Duration(0)

	renderFrame = js.NewCallback(func(args []js.Value) {
		if !browserRuntime.isPaused {
			now := time.Now()
			app.OnRender(elapsedFpsTime, mutex)
			elapsedFpsTime = fpsDelay - time.Since(now)
			if elapsedFpsTime < 0 {
				elapsedFpsTime = 0
			}
			time.Sleep(elapsedFpsTime)
		} else {
			time.Sleep(fpsDelay)
		}
		if !browserRuntime.isStopped {
			js.Global().Call("requestAnimationFrame", renderFrame)
		} else {
			browserRuntime.done <- true
		}
	})
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-browserRuntime.done

	renderFrame.Release()
	jsTge.Call("stop")

	return nil
}
