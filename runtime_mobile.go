// +build android ios

package tge

import (
	fmt "fmt"
	ioutil "io/ioutil"
	sync "sync"
	time "time"

	mobile "github.com/thommil/tge-mobile/app"
	asset "github.com/thommil/tge-mobile/asset"
	lifecycle "github.com/thommil/tge-mobile/event/lifecycle"
	paint "github.com/thommil/tge-mobile/event/paint"
	size "github.com/thommil/tge-mobile/event/size"
	touch "github.com/thommil/tge-mobile/event/touch"
	gl "github.com/thommil/tge-mobile/gl"
)

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type mobileRuntime struct {
	app            App
	plugins        map[string]Plugin
	host           mobile.App
	isPaused       bool
	isStopped      bool
	context        gl.Context
	lastMouseEvent MouseEvent
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

func (runtime *mobileRuntime) LoadAsset(p string) ([]byte, error) {
	if file, err := asset.Open(p); err != nil {
		return nil, err
	} else {
		return ioutil.ReadAll(file)
	}
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
		app:            app,
		plugins:        make(map[string]Plugin),
		isPaused:       true,
		isStopped:      true,
		lastMouseEvent: MouseEvent{},
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
					if mobileRuntime.context != nil && !e.External {
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
				switch e.Type {
				case touch.TypeBegin:
					// mouse down
					if (settings.EventMask & MouseButtonEventEnabled) != 0 {
						mobileRuntime.lastMouseEvent.X = int32(e.X)
						mobileRuntime.lastMouseEvent.Y = int32(e.Y)
						app.OnMouseEvent(
							MouseEvent{
								X:      mobileRuntime.lastMouseEvent.X,
								Y:      mobileRuntime.lastMouseEvent.Y,
								Type:   TypeDown,
								Button: ButtonNone,
							})
					}
				case touch.TypeMove:
					// mouse move
					if (settings.EventMask & MouseMotionEventEnabled) != 0 {
						x := int32(e.X)
						y := int32(e.Y)
						if (mobileRuntime.lastMouseEvent.X != x) && (mobileRuntime.lastMouseEvent.Y != y) {
							mobileRuntime.lastMouseEvent = MouseEvent{
								X:      x,
								Y:      y,
								Type:   TypeMove,
								Button: ButtonNone,
							}
							app.OnMouseEvent(mobileRuntime.lastMouseEvent)
						}
					}
				case touch.TypeEnd:
					// Touch down
					if (settings.EventMask & MouseButtonEventEnabled) != 0 {
						app.OnMouseEvent(
							MouseEvent{
								X:      int32(e.X),
								Y:      int32(e.Y),
								Type:   TypeUp,
								Button: ButtonNone,
							})
					}
				}
			}

		}
	})

	return nil
}
