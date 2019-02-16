// +build android ios

package tge

import (
	log "log"
	sync "sync"
	time "time"

	mobile "golang.org/x/mobile/app"
	lifecycle "golang.org/x/mobile/event/lifecycle"
	paint "golang.org/x/mobile/event/paint"
	size "golang.org/x/mobile/event/size"
	touch "golang.org/x/mobile/event/touch"
	gl "golang.org/x/mobile/gl"
)

type MobileRuntime interface {
	Runtime
	GetGlContext() gl.Context
}

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type mobileRuntime struct {
	app       App
	plugins   []Plugin
	mobile    mobile.App
	isPaused  bool
	ticker    *time.Ticker
	tickEnd   chan bool
	glContext gl.Context
}

func (runtime mobileRuntime) Use(plugin Plugin) {
	runtime.plugins = append(runtime.plugins, plugin)
	err := plugin.Init(runtime)
	if err != nil {
		log.Fatalln(err)
	}
}

func (runtime mobileRuntime) Stop() {
	// Not implemented
}

func (runtime mobileRuntime) GetGlContext() gl.Context {
	return runtime.glContext
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

	// Instanciate Runtime
	mobileRuntime := &mobileRuntime{
		app:      app,
		plugins:  make([]Plugin, 0),
		isPaused: true,
		tickEnd:  make(chan bool),
	}

	// -------------------------------------------------------------------- //
	// Ticker Loop
	// -------------------------------------------------------------------- //
	mutex := &sync.Mutex{}
	startTicker := func() {
		tpsDelay := time.Duration(1000000000 / settings.TPS)
		mobileRuntime.ticker = time.NewTicker(tpsDelay)

		elapsedTpsTime := time.Duration(0)
		go func() {
			for {
				select {
				case <-mobileRuntime.tickEnd:
					return
				case now := <-mobileRuntime.ticker.C:
					if !mobileRuntime.isPaused {
						app.OnTick(elapsedTpsTime, mutex)
						elapsedTpsTime = tpsDelay - time.Since(now)
						if elapsedTpsTime < 0 {
							elapsedTpsTime = 0
						}
					}
				}
			}
		}()
	}

	// -------------------------------------------------------------------- //
	// Init
	// -------------------------------------------------------------------- //
	fpsDelay := time.Duration(1000000000 / settings.FPS)
	elapsedFpsTime := time.Duration(0)
	mobile.Main(func(a mobile.App) {
		defer app.OnDispose()
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.To {
				case lifecycle.StageFocused:
					mobileRuntime.glContext, _ = e.DrawContext.(gl.Context)
					app.OnStart(mobileRuntime)
					startTicker()
					app.OnResume()
					mobileRuntime.isPaused = false

				case lifecycle.StageAlive:
					mobileRuntime.isPaused = true
					mobileRuntime.ticker.Stop()
					mobileRuntime.tickEnd <- true
					app.OnPause()
					app.OnStop()
					mobileRuntime.glContext = nil
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
				log.Println("OnTouch")

			}

		}
	})
	defer app.OnDispose()

	return nil
}
