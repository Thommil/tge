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
	key "golang.org/x/mobile/event/key"
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

// -------------------------------------------------------------------- //
// KeyMap
// -------------------------------------------------------------------- //

var keyMap = map[key.Code]KeyCode{
    key.CodeUnknown Code = 0

    key.CodeA : KeyCodeUnknown,
	key.CodeB : KeyCodeUnknown,    
	key.CodeC : KeyCodeUnknown,    
	key.CodeD : KeyCodeUnknown,    
	key.CodeE : KeyCodeUnknown,    
	key.CodeF : KeyCodeUnknown,    
	key.CodeG : KeyCodeUnknown,
    key.CodeH : KeyCodeUnknown,
    key.CodeI : KeyCodeUnknown,
    key.CodeJ : KeyCodeUnknown,
    key.CodeK : KeyCodeUnknown,
    key.CodeL : KeyCodeUnknown,
    key.CodeM : KeyCodeUnknown,
    key.CodeN : KeyCodeUnknown,
    key.CodeO : KeyCodeUnknown,
    key.CodeP : KeyCodeUnknown,
    key.CodeQ : KeyCodeUnknown,
    key.CodeR : KeyCodeUnknown,
    key.CodeS : KeyCodeUnknown,
    key.CodeT : KeyCodeUnknown,
    key.CodeU : KeyCodeUnknown,
    key.CodeV : KeyCodeUnknown,
    key.CodeW : KeyCodeUnknown,
    key.CodeX : KeyCodeUnknown,
    key.CodeY : KeyCodeUnknown,
    key.CodeZ : KeyCodeUnknown,

    key.Code1 : KeyCodeUnknown,
    key.Code2 : KeyCodeUnknown,
    key.Code3 : KeyCodeUnknown,
    key.Code4 : KeyCodeUnknown,
    key.Code5 : KeyCodeUnknown,
    key.Code6 : KeyCodeUnknown,
    key.Code7 : KeyCodeUnknown,
    key.Code8 : KeyCodeUnknown,
    key.Code9 : KeyCodeUnknown,
    key.Code0 : KeyCodeUnknown,

    key.CodeReturnEnter        : KeyCodeUnknown,
    key.CodeEscape             : KeyCodeUnknown,
    key.CodeDeleteBackspace    : KeyCodeUnknown,
    key.CodeTab                : KeyCodeUnknown,
    key.CodeSpacebar           : KeyCodeUnknown,
    key.CodeHyphenMinus        : KeyCodeUnknown, // -
    key.CodeEqualSign          : KeyCodeUnknown, // =
    key.CodeLeftSquareBracket  : KeyCodeUnknown, // [
    key.CodeRightSquareBracket : KeyCodeUnknown, // ]
    key.CodeBackslash          : KeyCodeUnknown, // \
    key.CodeSemicolon          : KeyCodeUnknown, // ;
    key.CodeApostrophe         : KeyCodeUnknown, // '
    key.CodeGraveAccent        : KeyCodeUnknown, // `
    key.CodeComma              : KeyCodeUnknown, // ,
    key.CodeFullStop           : KeyCodeUnknown, // .
    key.CodeSlash              : KeyCodeUnknown, // /
    key.CodeCapsLock           : KeyCodeUnknown,

    key.CodeF1  : KeyCodeUnknown,
    key.CodeF2  : KeyCodeUnknown,
    key.CodeF3  : KeyCodeUnknown,
    key.CodeF4  : KeyCodeUnknown,
    key.CodeF5  : KeyCodeUnknown,
    key.CodeF6  : KeyCodeUnknown,
    key.CodeF7  : KeyCodeUnknown,
    key.CodeF8  : KeyCodeUnknown,
    key.CodeF9  : KeyCodeUnknown,
    key.CodeF10 : KeyCodeUnknown,
    key.CodeF11 : KeyCodeUnknown,
    key.CodeF12 : KeyCodeUnknown,

    key.CodePause         : KeyCodeUnknown,
    key.CodeInsert        : KeyCodeUnknown,
    key.CodeHome          : KeyCodeUnknown,
    key.CodePageUp        : KeyCodeUnknown,
    key.CodeDeleteForward : KeyCodeUnknown,
    key.CodeEnd           : KeyCodeUnknown,
    key.CodePageDown      : KeyCodeUnknown,

    key.CodeRightArrow : KeyCodeUnknown,
    key.CodeLeftArrow  : KeyCodeUnknown,
    key.CodeDownArrow  : KeyCodeUnknown,
    key.CodeUpArrow    : KeyCodeUnknown,

    key.CodeKeypadNumLock     : KeyCodeUnknown,
    key.CodeKeypadSlash       : KeyCodeUnknown, // /
    key.CodeKeypadAsterisk    : KeyCodeUnknown, // *
    key.CodeKeypadHyphenMinus : KeyCodeUnknown, // -
    key.CodeKeypadPlusSign    : KeyCodeUnknown, // +
    key.CodeKeypadEnter       : KeyCodeUnknown,
    key.CodeKeypad1           : KeyCodeUnknown,
    key.CodeKeypad2           : KeyCodeUnknown,
    key.CodeKeypad3           : KeyCodeUnknown,
    key.CodeKeypad4           : KeyCodeUnknown,
    key.CodeKeypad5           : KeyCodeUnknown,
    key.CodeKeypad6           : KeyCodeUnknown,
    key.CodeKeypad7           : KeyCodeUnknown,
    key.CodeKeypad8           : KeyCodeUnknown,
    key.CodeKeypad9           : KeyCodeUnknown,
    key.CodeKeypad0           : KeyCodeUnknown,
    key.CodeKeypadFullStop    : KeyCodeUnknown,  // .
    key.CodeKeypadEqualSign   : KeyCodeUnknown, // =

    key.CodeHelp : KeyCodeUnknown,

    key.CodeMute       : KeyCodeUnknown,
    key.CodeVolumeUp   : KeyCodeUnknown,
    key.CodeVolumeDown : KeyCodeUnknown,

    key.CodeLeftControl  : KeyCodeUnknown,
    key.CodeLeftShift    : KeyCodeUnknown,
    key.CodeLeftAlt      : KeyCodeUnknown,
    key.CodeLeftGUI      : KeyCodeUnknown,
    key.CodeRightControl : KeyCodeUnknown,
    key.CodeRightShift   : KeyCodeUnknown,
    key.CodeRightAlt     : KeyCodeUnknown,
    key.CodeRightGUI     : KeyCodeUnknown,

    // CodeCompose is the Code for a compose key, sometimes called a multi key,
    // used to input non-ASCII characters such as ñ being composed of n and ~.
    //
    // See https://en.wikipedia.org/wiki/Compose_key
    key.CodeCompose : KeyCodeUnknown,
)
