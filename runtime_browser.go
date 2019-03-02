// Copyright (c) 2019 Thomas MILLET. All rights reserved.

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
	ticker    *time.Ticker
	canvas    *js.Value
	jsTge     *js.Value
	isPaused  bool
	isStopped bool
	done      chan bool
}

func (runtime *browserRuntime) Use(plugin Plugin) {
	use(plugin, runtime)
}

func (runtime *browserRuntime) GetAsset(p string) ([]byte, error) {
	var data []byte
	var err error
	var doneState = make(chan bool)

	onLoadAssetCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0] != js.Null() {
			err = fmt.Errorf(args[1].String())
			doneState <- false
		} else {
			doneState <- true
		}
		return false
	})

	onGetAssetSizeCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[1] != js.Null() {
			err = fmt.Errorf(args[1].String())
			doneState <- false
		} else if size := args[0].Int(); size > 0 {
			data = make([]byte, size)
			jsData := js.TypedArrayOf(data)
			defer jsData.Release()
			runtime.jsTge.Call("loadAsset", p, jsData, onLoadAssetCallback)
		} else {
			err = fmt.Errorf("empty asset")
			doneState <- false
		}
		return false
	})
	defer onLoadAssetCallback.Release()
	defer onGetAssetSizeCallback.Release()

	runtime.jsTge.Call("getAssetSize", p, onGetAssetSizeCallback)

	<-doneState

	return data, err
}

func (runtime *browserRuntime) GetHost() interface{} {
	host := js.Global()
	return &host
}

func (runtime *browserRuntime) GetPlugin(name string) Plugin {
	return plugins[name]
}

func (runtime *browserRuntime) GetRenderer() interface{} {
	glContext := runtime.canvas.Call("getContext", "webgl2")
	if glContext == js.Undefined() {
		fmt.Println("WARNING: No WebGL2 support")
		glContext = runtime.canvas.Call("getContext", "webgl")
	}
	if glContext == js.Undefined() {
		fmt.Println("WARNING: No WebGL support")
		glContext = runtime.canvas.Call("getContext", "experimental-webgl")
	}
	if glContext == js.Undefined() {
		err := fmt.Errorf("No WebGL support found in brower")
		fmt.Println(err)
		panic(err)
	}
	return &glContext
}

func (runtime *browserRuntime) Subscribe(channel string, listener Listener) {
	subscribe(channel, listener)
}

func (runtime *browserRuntime) Unsubscribe(channel string, listener Listener) {
	unsubscribe(channel, listener)
}

func (runtime *browserRuntime) Publish(event Event) {
	publish(event)
}

func (runtime *browserRuntime) Stop() {
	if !runtime.isPaused {
		runtime.isPaused = true
		runtime.app.OnPause()
	}
	runtime.isStopped = true
	runtime.app.OnStop()
	dispose()
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
		isPaused:  true,
		isStopped: true,
		canvas:    &canvas,
		done:      make(chan bool),
		jsTge:     &jsTge,
	}

	// Start App
	err = app.OnStart(browserRuntime)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	browserRuntime.isStopped = false

	// Resume App
	app.OnResume()
	browserRuntime.isPaused = false

	// Resize App
	go publish(ResizeEvent{int32(browserRuntime.canvas.Get("clientWidth").Int()),
		int32(browserRuntime.canvas.Get("clientHeight").Int())})

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
	resizeEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !browserRuntime.isStopped {
			publish(ResizeEvent{int32(browserRuntime.canvas.Get("clientWidth").Int()),
				int32(browserRuntime.canvas.Get("clientHeight").Int())})
		}
		return false
	})
	defer resizeEvtCb.Release()
	js.Global().Call("addEventListener", "resize", resizeEvtCb)

	// Focus
	blurEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !browserRuntime.isStopped && !browserRuntime.isPaused {
			browserRuntime.isPaused = true
			browserRuntime.app.OnPause()
		}
		return false
	})
	defer blurEvtCb.Release()
	browserRuntime.canvas.Call("addEventListener", "blur", blurEvtCb)

	focuseEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !browserRuntime.isStopped && browserRuntime.isPaused {
			//Called in go routine in case of asset loading in resume (blocking)
			go func() {
				browserRuntime.app.OnResume()
				browserRuntime.isPaused = false
			}()
		}
		return false
	})
	defer focuseEvtCb.Release()
	browserRuntime.canvas.Call("addEventListener", "focus", focuseEvtCb)

	// Destroy
	beforeunloadEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !browserRuntime.isStopped {
			browserRuntime.Stop()
		}
		return false
	})
	defer beforeunloadEvtCb.Release()
	js.Global().Call("addEventListener", "beforeunload", beforeunloadEvtCb)

	// MouseButtonEvent
	if (settings.EventMask & MouseButtonEventEnabled) != 0 {
		mouseDownEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if !browserRuntime.isStopped && !browserRuntime.isPaused {
				event := args[0]
				publish(MouseEvent{
					X:      int32(event.Get("offsetX").Int()),
					Y:      int32(event.Get("offsetY").Int()),
					Button: Button(event.Get("button").Int() + 1),
					Type:   TypeDown,
				})
			}
			return false
		})
		defer mouseDownEvtCb.Release()
		browserRuntime.canvas.Call("addEventListener", "mousedown", mouseDownEvtCb)

		mouseUpEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if !browserRuntime.isStopped && !browserRuntime.isPaused {
				event := args[0]
				publish(MouseEvent{
					X:      int32(event.Get("offsetX").Int()),
					Y:      int32(event.Get("offsetY").Int()),
					Button: Button(event.Get("button").Int() + 1),
					Type:   TypeUp,
				})
			}
			return false
		})
		defer mouseUpEvtCb.Release()
		browserRuntime.canvas.Call("addEventListener", "mouseup", mouseUpEvtCb)
	}

	// MouseMotionEventEnabled
	if (settings.EventMask & MouseMotionEventEnabled) != 0 {
		mouseMoveEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if !browserRuntime.isStopped && !browserRuntime.isPaused {
				event := args[0]
				offsetX := event.Get("offsetX").Int()
				offsetY := event.Get("offsetY").Int()
				if (offsetX > settings.MouseMotionThreshold) || (offsetY > settings.MouseMotionThreshold) {
					publish(MouseEvent{
						X:      int32(offsetX),
						Y:      int32(offsetY),
						Button: ButtonNone,
						Type:   TypeMove,
					})
				}
			}
			return false
		})
		defer mouseMoveEvtCb.Release()
		browserRuntime.canvas.Call("addEventListener", "mousemove", mouseMoveEvtCb)
	}

	// ScrollEvent
	if (settings.EventMask & ScrollEventEnabled) != 0 {
		wheelEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if !browserRuntime.isStopped && !browserRuntime.isPaused {
				event := args[0]
				publish(ScrollEvent{
					X: int32(event.Get("deltaX").Int()),
					Y: -int32(event.Get("deltaY").Int()),
				})
			}
			return false
		})
		defer wheelEvtCb.Release()
		browserRuntime.canvas.Call("addEventListener", "wheel", wheelEvtCb)
	}

	// KeyEvent
	if (settings.EventMask & KeyEventEnabled) != 0 {
		keyDownEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if !browserRuntime.isStopped && !browserRuntime.isPaused {
				event := args[0]
				event.Call("preventDefault")
				keyCode := event.Get("key").String()
				publish(KeyEvent{
					Key:   keyMap[keyCode],
					Value: keyCode,
					Type:  TypeDown,
				})
			}
			return false
		})
		defer keyDownEvtCb.Release()
		browserRuntime.canvas.Call("addEventListener", "keydown", keyDownEvtCb)

		keyUpEvtCb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if !browserRuntime.isStopped && !browserRuntime.isPaused {
				event := args[0]
				event.Call("preventDefault")
				keyCode := event.Get("key").String()
				publish(KeyEvent{
					Key:   keyMap[keyCode],
					Value: keyCode,
					Type:  TypeUp,
				})
			}
			return false
		})
		defer keyUpEvtCb.Release()
		browserRuntime.canvas.Call("addEventListener", "keyup", keyUpEvtCb)
	}

	// -------------------------------------------------------------------- //
	// Render Loop
	// -------------------------------------------------------------------- //
	var renderFrame js.Func
	elapsedFpsTime := time.Duration(0)

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !browserRuntime.isPaused {
			now := time.Now()
			app.OnRender(elapsedFpsTime, mutex)
			elapsedFpsTime = time.Since(now)
		}
		if !browserRuntime.isStopped {
			js.Global().Call("requestAnimationFrame", renderFrame)
		} else {
			browserRuntime.done <- true
		}
		return false
	})
	js.Global().Call("requestAnimationFrame", renderFrame)

	<-browserRuntime.done

	renderFrame.Release()
	jsTge.Call("stop")

	noExit := make(chan int)
	<-noExit

	return nil
}

// -------------------------------------------------------------------- //
// KeyMap
// -------------------------------------------------------------------- //

var keyMap = map[string]KeyCode{

	// Printable
	"a": KeyCodeA,
	"b": KeyCodeB,
	"c": KeyCodeC,
	"d": KeyCodeD,
	"e": KeyCodeE,
	"f": KeyCodeF,
	"g": KeyCodeG,
	"h": KeyCodeH,
	"i": KeyCodeI,
	"j": KeyCodeJ,
	"k": KeyCodeK,
	"l": KeyCodeL,
	"m": KeyCodeM,
	"n": KeyCodeN,
	"o": KeyCodeO,
	"p": KeyCodeP,
	"q": KeyCodeQ,
	"r": KeyCodeR,
	"s": KeyCodeS,
	"t": KeyCodeT,
	"u": KeyCodeU,
	"v": KeyCodeV,
	"w": KeyCodeW,
	"x": KeyCodeX,
	"y": KeyCodeY,
	"z": KeyCodeZ,

	"1": KeyCode1,
	"2": KeyCode2,
	"3": KeyCode3,
	"4": KeyCode4,
	"5": KeyCode5,
	"6": KeyCode6,
	"7": KeyCode7,
	"8": KeyCode8,
	"9": KeyCode9,
	"0": KeyCode0,

	"Enter": KeyCodeReturnEnter,
	"Tab":   KeyCodeTab,
	" ":     KeyCodeSpacebar,
	"-":     KeyCodeHyphenMinus,        // -
	"=":     KeyCodeEqualSign,          // =
	"[":     KeyCodeLeftSquareBracket,  // [
	"]":     KeyCodeRightSquareBracket, // ]
	"\\":    KeyCodeBackslash,          // \
	";":     KeyCodeSemicolon,          // ;
	"'":     KeyCodeApostrophe,         // '
	"`":     KeyCodeGraveAccent,        // `
	",":     KeyCodeComma,              // ,
	".":     KeyCodeFullStop,           // .
	"/":     KeyCodeSlash,              // /

	"Divide":    KeyCodeKeypadSlash,       // /
	"Multiply":  KeyCodeKeypadAsterisk,    // *
	"Substract": KeyCodeKeypadHyphenMinus, // -
	"Add":       KeyCodeKeypadPlusSign,    // +
	"Decimal":   KeyCodeKeypadFullStop,    // .

	"@": KeyCodeAt,                // @
	">": KeyCodeGreaterThan,       // >
	"<": KeyCodeLesserThan,        // <
	"$": KeyCodeDollar,            // $
	":": KeyCodeColon,             // :
	"(": KeyCodeLeftParenthesis,   // (
	")": KeyCodeLRightParenthesis, // )

	"&":  KeyCodeAmpersand,   // &
	"#":  KeyCodeHash,        // #
	"\"": KeyDoubleQuote,     // "
	"''": KeyQuote,           // '
	"§":  KeyParapgrah,       // §
	"!":  KeyExclamationMark, // !
	"_":  KeyUnderscore,      // _
	"?":  KeyQuestionMark,    // ?
	"%":  KeyPercent,         // %
	"°":  KeyDegree,          // °

	// Actions

	"Escape":   KeyCodeEscape,
	"CapsLock": KeyCodeCapsLock,

	"Backspace": KeyCodeDeleteBackspace,
	"Pause":     KeyCodePause,
	"Insert":    KeyCodeInsert,
	"Home":      KeyCodeHome,
	"PageUp":    KeyCodePageUp,
	"Delete":    KeyCodeDeleteForward,
	"End":       KeyCodeEnd,
	"PageDown":  KeyCodePageDown,

	"ArrowRight": KeyCodeRightArrow,
	"ArrowLeft":  KeyCodeLeftArrow,
	"ArrowDown":  KeyCodeDownArrow,
	"ArrowUp":    KeyCodeUpArrow,

	"Numlock": KeyCodeKeypadNumLock,

	"Help": KeyCodeHelp,

	"AudioVolumeMute": KeyCodeMute,
	"AudioVolumeUp":   KeyCodeVolumeUp,
	"AudioVolumeDown": KeyCodeVolumeDown,

	// Functions

	"F1":  KeyCodeF1,
	"F2":  KeyCodeF2,
	"F3":  KeyCodeF3,
	"F4":  KeyCodeF4,
	"F5":  KeyCodeF5,
	"F6":  KeyCodeF6,
	"F7":  KeyCodeF7,
	"F8":  KeyCodeF8,
	"F9":  KeyCodeF9,
	"F10": KeyCodeF10,
	"F11": KeyCodeF11,
	"F12": KeyCodeF12,

	// Modifiers

	"Control": KeyCodeLeftControl,
	"Shift":   KeyCodeLeftShift,
	"Alt":     KeyCodeLeftAlt,
	"Meta":    KeyCodeLeftGUI,

	// Compose

	"Compose": KeyCodeCompose,
}
