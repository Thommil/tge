// +build darwin freebsd linux windows
// +build !android
// +build !ios

package tge

import (
	log "log"
	runtime "runtime"

	glfw "github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

func doRun(app App, settings *Settings) error {
	log.Println("doRun()")
	err := glfw.Init()
	if err != nil {
		return err
	}
	defer glfw.Terminate()

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

	window, err := glfw.CreateWindow(settings.Width, settings.Height, settings.Name, monitor, nil)
	if err != nil {
		return err
	}

	window.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		app.Resize(width, height)
	})

	window.SetCloseCallback(func(w *glfw.Window) {
		app.Pause()
		app.Stop()
	})

	window.SetFocusCallback(func(w *glfw.Window, focused bool) {
		if focused {
			app.Resume()
		} else {
			app.Pause()
		}
	})

	window.MakeContextCurrent()
	app.Start()
	window.Focus()

	for !window.ShouldClose() {
		//app.Render()
		//app.Tick()
		window.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}
