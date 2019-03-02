// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// +build android ios

package tge

import (
	fmt "fmt"
	ioutil "io/ioutil"
	"math"
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
	host           mobile.App
	isPaused       bool
	isStopped      bool
	context        gl.Context
	lastMouseEvent MouseEvent
}

func (runtime *mobileRuntime) Use(plugin Plugin) {
	use(plugin, runtime)
}

func (runtime *mobileRuntime) GetAsset(p string) ([]byte, error) {
	if file, err := asset.Open(p); err != nil {
		return nil, err
	} else {
		return ioutil.ReadAll(file)
	}
}

func (runtime *mobileRuntime) GetHost() interface{} {
	return runtime.host
}

func (runtime *mobileRuntime) GetPlugin(name string) Plugin {
	return plugins[name]
}

func (runtime *mobileRuntime) GetRenderer() interface{} {
	return runtime.context
}

func (runtime *mobileRuntime) Subscribe(channel string, listener Listener) {
	subscribe(channel, listener)
}

func (runtime *mobileRuntime) Unsubscribe(channel string, listener Listener) {
	unsubscribe(channel, listener)
}

func (runtime *mobileRuntime) Publish(event Event) {
	publish(event)
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
		isPaused:       true,
		isStopped:      true,
		lastMouseEvent: MouseEvent{},
	}
	defer dispose()

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
	var moveEvtChan chan MouseEvent
	elapsedFpsTime := time.Duration(0)
	mobile.Main(func(a mobile.App) {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.To {
				case lifecycle.StageFocused:
					mobileRuntime.context, _ = e.DrawContext.(gl.Context)
					mobileRuntime.host = a
					err := app.OnStart(mobileRuntime)
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
					mobileRuntime.isStopped = false
					go startTicker()
					app.OnResume()
					mobileRuntime.isPaused = false

					// Mouse motion hack to queue move events
					moveEvtChan = make(chan MouseEvent, 100)
					go func() {
						for !mobileRuntime.isStopped {
							publish(<-moveEvtChan)
						}
					}()

				case lifecycle.StageAlive:
					mobileRuntime.isPaused = true
					app.OnPause()
					mobileRuntime.isStopped = true
					close(moveEvtChan)
					app.OnStop()
					mobileRuntime.context = nil
				}

			case paint.Event:
				if !mobileRuntime.isPaused {
					if mobileRuntime.context != nil && !e.External {
						now := time.Now()
						app.OnRender(elapsedFpsTime, mutex)
						a.Publish()
						elapsedFpsTime = time.Since(now)
					}
					a.Send(paint.Event{})
				}

			case size.Event:
				go publish(ResizeEvent{int32(e.WidthPx), int32(e.HeightPx)})

			case touch.Event:
				switch e.Type {
				case touch.TypeBegin:
					// mouse down
					if (settings.EventMask & MouseButtonEventEnabled) != 0 {
						mobileRuntime.lastMouseEvent.X = int32(e.X)
						mobileRuntime.lastMouseEvent.Y = int32(e.Y)
						go publish(MouseEvent{
							X:      mobileRuntime.lastMouseEvent.X,
							Y:      mobileRuntime.lastMouseEvent.Y,
							Type:   TypeDown,
							Button: ButtonNone,
						})
					}
				case touch.TypeMove:
					// mouse move
					if (settings.EventMask & MouseMotionEventEnabled) != 0 {
						go func() {
							x := int32(e.X)
							y := int32(e.Y)
							if math.Abs(float64(mobileRuntime.lastMouseEvent.X-x)) > float64(settings.MouseMotionThreshold) || math.Abs(float64(mobileRuntime.lastMouseEvent.Y-y)) > float64(settings.MouseMotionThreshold) {
								moveEvtChan <- MouseEvent{
									X:      x,
									Y:      y,
									Type:   TypeMove,
									Button: ButtonNone,
								}
							}
						}()
					}
				case touch.TypeEnd:
					// Touch down
					if (settings.EventMask & MouseButtonEventEnabled) != 0 {
						mobileRuntime.lastMouseEvent.X = 0
						mobileRuntime.lastMouseEvent.Y = 0
						go publish(MouseEvent{
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
