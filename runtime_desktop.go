// +build darwin freebsd linux windows
// +build !android
// +build !ios

package tge

import (
	log "log"
	runtime "runtime"
	sync "sync"

	glfw "github.com/go-gl/glfw/v3.2/glfw"
	physics "github.com/thommil/tge/physics"
	player "github.com/thommil/tge/player"
	renderer "github.com/thommil/tge/renderer"
	ui "github.com/thommil/tge/ui"
)

// Runtime implementation
type desktopRuntime struct {
	window   *glfw.Window
	renderer renderer.Renderer
	ui       ui.UI
	player   player.Player
	physics  physics.Physics
}

func (runtime desktopRuntime) Stop() {
	runtime.window.SetShouldClose(true)
}

// init ensure that we're running on main thread
func init() {
	runtime.LockOSThread()
}

// Main
func doRun(app App, settings *Settings) error {
	log.Println("doRun()")

	// Init
	err := glfw.Init()
	if err != nil {
		return err
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
		return err
	}

	// Window callbacks
	var resizeAtStart sync.Once
	var paused = true
	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		// Windows minify issue
		if !paused && width > 0 {
			app.OnResize(width, height)
		}
	})

	window.SetCloseCallback(func(w *glfw.Window) {
		if !paused {
			app.OnPause()
		}
		app.OnStop()
	})

	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused && paused {
			paused = false
			app.OnResume()
			resizeAtStart.Do(func() {
				if runtime.GOOS != "windows" {
					app.OnResize(settings.Width, settings.Height)
				}
			})
		} else if !paused {
			paused = true
			app.OnPause()
		}
	})

	// Start
	window.MakeContextCurrent()

	desktopRuntime := desktopRuntime{
		window: window,
	}
	app.OnStart(&desktopRuntime)

	// OS Specific
	if runtime.GOOS == "windows" {
		paused = false
		app.OnResume()
		app.OnResize(settings.Width, settings.Height)
	}

	// Loop
	for !window.ShouldClose() {
		//app.Render()
		//app.Tick()
		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}
