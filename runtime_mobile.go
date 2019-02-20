// +build android ios

package tge

import (
	fmt "fmt"
	sync "sync"
	time "time"

	mobile "golang.org/x/mobile/app"
	lifecycle "golang.org/x/mobile/event/lifecycle"
	paint "golang.org/x/mobile/event/paint"
	size "golang.org/x/mobile/event/size"
	touch "golang.org/x/mobile/event/touch"
	gl "golang.org/x/mobile/gl"
)

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type mobileRuntime struct {
	app       App
	plugins   map[string]Plugin
	host      mobile.App
	isPaused  bool
	isStopped bool
	context   gl.Context
}

func (runtime *mobileRuntime) Use(plugin Plugin) {
	name := plugin.GetName()
	if _, found := runtime.plugins[name]; !found {
		runtime.plugins[name] = plugin
		err := plugin.Init(runtime)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Printf("Plugin %s loaded\n", name)
	}
}

func (runtime *mobileRuntime) GetPlugin(name string) Plugin {
	return runtime.plugins[name]
}

func (runtime *mobileRuntime) GetRenderer() interface{} {
	return runtime.context
}

func (runtime *mobileRuntime) GetHost() interface{} {
	return runtime.host
}

func (runtime *mobileRuntime) Stop() {
	// Not implemented
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
	defer app.OnDispose()

	// Instanciate Runtime
	mobileRuntime := &mobileRuntime{
		app:       app,
		plugins:   make(map[string]Plugin),
		isPaused:  true,
		isStopped: true,
	}

	// Unload plugins
	defer func() {
		for _, plugin := range mobileRuntime.plugins {
			plugin.Dispose()
		}
	}()

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	mutex := &sync.Mutex{}
	startTicker := func() {
		tpsDelay := time.Duration(1000000000 / settings.TPS)
		elapsedTpsTime := time.Duration(0)
		for !mobileRuntime.isStopped {
			if !mobileRuntime.isPaused {
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
	}

	// -------------------------------------------------------------------- //
	// Init
	// -------------------------------------------------------------------- //
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	elapsedFpsTime := time.Duration(0)
	mobile.Main(func(a mobile.App) {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.To {
				case lifecycle.StageFocused:
					mobileRuntime.context, _ = e.DrawContext.(gl.Context)
					mobileRuntime.host = a
					app.OnStart(mobileRuntime)
					mobileRuntime.isStopped = false
					go startTicker()
					app.OnResume()
					mobileRuntime.isPaused = false

				case lifecycle.StageAlive:
					mobileRuntime.isPaused = true
					app.OnPause()
					mobileRuntime.isStopped = true
					app.OnStop()
					mobileRuntime.context = nil
				}

			case paint.Event:
				if !mobileRuntime.isPaused {
					if mobileRuntime.glContext != nil && !e.External {
						now := time.Now()
						app.OnRender(elapsedFpsTime, mutex)
						a.Publish()
						elapsedFpsTime = fpsDelay - time.Since(now)
						if elapsedFpsTime < 0 {
							elapsedFpsTime = 0
						}
						time.Sleep(elapsedFpsTime)
					}
					a.Send(paint.Event{})
				}

			case size.Event:

				app.OnResize(e.WidthPx, e.HeightPx)

			case touch.Event:
				fmt.Println("OnTouch")

			}

		}
	})

	return nil
}
