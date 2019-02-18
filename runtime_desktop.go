// +build darwin freebsd linux windows
// +build !android
// +build !ios
// +build !js

package tge

import (
	fmt "fmt"
	runtime "runtime"
	sync "sync"
	time "time"

	sdl "github.com/veandco/go-sdl2/sdl"
)

// init ensure that we're running on main thread
func init() {
	runtime.LockOSThread()
}

type DesktopRuntime interface {
	Runtime
	GetWindow() *sdl.Window
}

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type desktopRuntime struct {
	app       App
	plugins   map[string]Plugin
	window    *sdl.Window
	isPaused  bool
	isStopped bool
}

func (runtime *desktopRuntime) Use(plugin Plugin) {
	runtime.plugins[plugin.GetName()] = plugin
	err := plugin.Init(runtime)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (runtime *desktopRuntime) GetPlugin(name string) Plugin {
	return runtime.plugins[name]
}

func (runtime *desktopRuntime) GetWindow() *sdl.Window {
	return runtime.window
}

func (runtime *desktopRuntime) Stop() {
	if !runtime.isPaused {
		runtime.isPaused = true
		runtime.app.OnPause()
	}
	runtime.isStopped = true
	runtime.app.OnStop()
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
	defer app.OnDispose()

	// -------------------------------------------------------------------- //
	// Init
	// -------------------------------------------------------------------- //
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer sdl.Quit()

	windowFlags := sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE
	if settings.Fullscreen {
		windowFlags = windowFlags | sdl.WINDOW_FULLSCREEN_DESKTOP
	}

	// Window creation
	window, err := sdl.CreateWindow(settings.Name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(settings.Width), int32(settings.Height), uint32(windowFlags))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer window.Destroy()

	context, err := window.GLCreateContext()
	if err != nil {
		panic(err)
	}

	// Instanciate Runtime
	desktopRuntime := &desktopRuntime{
		app:       app,
		plugins:   make(map[string]Plugin),
		window:    window,
		isPaused:  true,
		isStopped: true,
	}

	// Unload plugins
	defer func() {
		for _, plugin := range desktopRuntime.plugins {
			plugin.Dispose()
		}
	}()
	defer sdl.GLDeleteContext(context)

	// Start App
	app.OnStart(desktopRuntime)
	desktopRuntime.isStopped = false

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	mutex := &sync.Mutex{}
	tpsDelay := time.Duration(1000000000 / settings.TPS)
	elapsedTpsTime := time.Duration(0)
	go func() {
		for !desktopRuntime.isStopped {
			if !desktopRuntime.isPaused {
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
	// Render Loop
	// -------------------------------------------------------------------- //
	var resizeAtStart sync.Once
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	elapsedFpsTime := time.Duration(0)
	for !desktopRuntime.isStopped {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				desktopRuntime.Stop()
			case *sdl.WindowEvent:
				switch t.Event {
				case sdl.WINDOWEVENT_FOCUS_GAINED:
					app.OnResume()
					desktopRuntime.isPaused = false
					resizeAtStart.Do(func() {
						w, h := window.GetSize()
						app.OnResize(int(w), int(h))
					})
				case sdl.WINDOWEVENT_FOCUS_LOST:
					desktopRuntime.isPaused = true
					app.OnPause()
				case sdl.WINDOWEVENT_RESIZED:
					w, h := window.GetSize()
					app.OnResize(int(w), int(h))
				}
			case *sdl.MouseMotionEvent:

			}
		}
		if !desktopRuntime.isPaused {
			now := time.Now()
			app.OnRender(elapsedFpsTime, mutex)
			window.GLSwap()
			elapsedFpsTime = fpsDelay - time.Since(now)
			if elapsedFpsTime < 0 {
				elapsedFpsTime = 0
			}
			time.Sleep(elapsedFpsTime)
		}
	}

	return nil
}
