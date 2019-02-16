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

	glfw "github.com/go-gl/glfw/v3.2/glfw"
)

// init ensure that we're running on main thread
func init() {
	runtime.LockOSThread()
}

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type desktopRuntime struct {
	app       App
	plugins   []Plugin
	window    *glfw.Window
	isPaused  bool
	isStopped bool
}

func (runtime *desktopRuntime) Use(plugin Plugin) {
	runtime.plugins = append(runtime.plugins, plugin)
	err := plugin.Init(runtime)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
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
	err = glfw.Init()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer glfw.Terminate()

	// Fullscreen support
	var monitor *glfw.Monitor
	if settings.Fullscreen {
		monitor = glfw.GetPrimaryMonitor()
		videoMode := monitor.GetVideoMode()
		settings.Width = videoMode.Width
		settings.Height = videoMode.Height
		glfw.WindowHint(glfw.RedBits, videoMode.RedBits)
		glfw.WindowHint(glfw.GreenBits, videoMode.GreenBits)
		glfw.WindowHint(glfw.BlueBits, videoMode.BlueBits)
		glfw.WindowHint(glfw.RefreshRate, videoMode.RefreshRate)
	}

	// Window creation
	window, err := glfw.CreateWindow(settings.Width, settings.Height, settings.Name, monitor, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer window.SetShouldClose(true)

	// Start GLFW
	window.MakeContextCurrent()

	// Instanciate Runtime
	desktopRuntime := &desktopRuntime{
		app:       app,
		plugins:   make([]Plugin, 0),
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

	// Start App
	app.OnStart(desktopRuntime)
	desktopRuntime.isStopped = false

	// OS Specific - Windows do not focus at start
	if runtime.GOOS == "windows" {
		app.OnResume()
		desktopRuntime.isPaused = false
		app.OnResize(settings.Width, settings.Height)
	}

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
	// Callbacks
	// -------------------------------------------------------------------- //
	var resizeAtStart sync.Once

	// Resize
	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		// OS Specific - Windows call resize to 0
		if !desktopRuntime.isPaused && width > 0 {
			app.OnResize(width, height)
		}
	})

	// Focus
	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused && desktopRuntime.isPaused {
			app.OnResume()
			desktopRuntime.isPaused = false
			// OS Specific - MacOS do not resize at start
			resizeAtStart.Do(func() {
				if runtime.GOOS != "windows" {
					app.OnResize(settings.Width, settings.Height)
				}
			})
		} else if !desktopRuntime.isPaused {
			desktopRuntime.isPaused = true
			app.OnPause()
		}
	})

	// Destroy
	window.SetCloseCallback(func(w *glfw.Window) {
		desktopRuntime.Stop()
	})

	// -------------------------------------------------------------------- //
	// Render Loop
	// -------------------------------------------------------------------- //
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	elapsedFpsTime := time.Duration(0)
	for !desktopRuntime.isStopped {
		if !desktopRuntime.isPaused {
			now := time.Now()
			app.OnRender(elapsedFpsTime, mutex)
			window.SwapBuffers()
			elapsedFpsTime = fpsDelay - time.Since(now)
			if elapsedFpsTime < 0 {
				elapsedFpsTime = 0
			}
			time.Sleep(elapsedFpsTime)
		}
		glfw.PollEvents()
	}

	return nil
}
