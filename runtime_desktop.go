// Copyright (c) 2019 Thomas MILLET. All rights reserved.

// Copyright (c) 2013, Go-SDL2 Authors
// All rights reserved.

// +build darwin freebsd linux windows
// +build !android
// +build !ios
// +build !js

package tge

import (
	fmt "fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	runtime "runtime"
	sync "sync"
	time "time"

	sdl "github.com/veandco/go-sdl2/sdl"
)

// init ensure that we're running on main thread
func init() {
	runtime.LockOSThread()
}

// -------------------------------------------------------------------- //
// Runtime implementation
// -------------------------------------------------------------------- //
type desktopRuntime struct {
	app        App
	host       *sdl.Window
	context    *sdl.GLContext
	isPaused   bool
	isStopped  bool
	assetsPath string
}

func (runtime *desktopRuntime) Use(plugin Plugin) {
	use(plugin, runtime)
}

func (runtime *desktopRuntime) GetAsset(p string) ([]byte, error) {
	return ioutil.ReadFile(path.Join(runtime.assetsPath, p))
}

func (runtime *desktopRuntime) GetHost() interface{} {
	return runtime.host
}

func (runtime *desktopRuntime) GetPlugin(name string) Plugin {
	return plugins[name]
}

func (runtime *desktopRuntime) GetRenderer() interface{} {
	return runtime.context
}

func (runtime *desktopRuntime) Subscribe(channel string, listener Listener) {
	subscribe(channel, listener)
}

func (runtime *desktopRuntime) Unsubscribe(channel string, listener Listener) {
	unsubscribe(channel, listener)
}

func (runtime *desktopRuntime) Publish(event Event) {
	publish(event)
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
	sdl.SetHint(sdl.HINT_VIDEO_HIGHDPI_DISABLED, "1")
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer sdl.Quit()

	windowFlags := sdl.WINDOW_OPENGL | sdl.WINDOW_RESIZABLE
	if settings.Fullscreen {
		windowFlags = windowFlags | sdl.WINDOW_FULLSCREEN_DESKTOP
	}

	if runtime.GOOS == "darwin" {
		sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	}

	sdl.GLSetAttribute(sdl.GL_MULTISAMPLEBUFFERS, 1)
	sdl.GLSetAttribute(sdl.GL_MULTISAMPLESAMPLES, 2)

	sdl.GLSetAttribute(sdl.GL_ACCELERATED_VISUAL, 1)

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
		host:      window,
		context:   &context,
		isPaused:  true,
		isStopped: true,
	}

	// Eval assets path
	if p, err := os.Executable(); err != nil {
		panic(err)
	} else {
		if p, err = filepath.EvalSymlinks(p); err != nil {
			panic(err)
		}
		if runtime.GOOS == "darwin" {
			// Packed mode (DIST for darwin)
			desktopRuntime.assetsPath = path.Join(path.Dir(p), "../Resources")
		} else {
			// Unpacked mode (DIST for windows/linux)
			desktopRuntime.assetsPath = path.Join(path.Dir(p), "assets")
		}

		if _, err := os.Stat(desktopRuntime.assetsPath); os.IsNotExist(err) {
			// Unpacked mode (DEV for all)
			desktopRuntime.assetsPath = path.Join(path.Dir(p), "../../assets")
		}

	}

	// Unload plugins
	defer dispose()
	defer sdl.GLDeleteContext(context)

	// Start App
	err = app.OnStart(desktopRuntime)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
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
						publish(ResizeEvent{w, h})
					})
				case sdl.WINDOWEVENT_FOCUS_LOST:
					desktopRuntime.isPaused = true
					app.OnPause()
				case sdl.WINDOWEVENT_RESIZED:
					w, h := window.GetSize()
					publish(ResizeEvent{w, h})
				}
			case *sdl.MouseButtonEvent:
				if (settings.EventMask & MouseButtonEventEnabled) != 0 {
					publish(MouseEvent{
						X:      t.X,
						Y:      t.Y,
						Type:   Type(t.Type),
						Button: Button(t.Button),
					})
				}
			case *sdl.MouseMotionEvent:
				if (settings.EventMask & MouseMotionEventEnabled) != 0 {
					if (int(t.X) > settings.MouseMotionThreshold) || (int(t.Y) > settings.MouseMotionThreshold) {
						publish(MouseEvent{
							X:      t.X,
							Y:      t.Y,
							Type:   TypeMove,
							Button: ButtonNone,
						})
					}
				}
			case *sdl.MouseWheelEvent:
				if (settings.EventMask & ScrollEventEnabled) != 0 {
					publish(ScrollEvent{
						X: t.X,
						Y: t.Y,
					})
				}
			case *sdl.KeyboardEvent:
				if (settings.EventMask & KeyEventEnabled) != 0 {
					keyCode := sdl.GetKeyName(t.Keysym.Sym)
					publish(KeyEvent{
						Type:  Type(t.Type),
						Key:   keyMap[keyCode],
						Value: keyCode,
					})
				}
			}
		}
		if !desktopRuntime.isPaused {
			now := time.Now()
			app.OnRender(elapsedFpsTime, mutex)
			window.GLSwap()
			elapsedFpsTime = time.Since(now)
		}
	}

	return nil
}

// -------------------------------------------------------------------- //
// KeyMap
// -------------------------------------------------------------------- //

var keyMap = map[string]KeyCode{

	// Printable
	"A": KeyCodeA,
	"B": KeyCodeB,
	"C": KeyCodeC,
	"D": KeyCodeD,
	"E": KeyCodeE,
	"F": KeyCodeF,
	"G": KeyCodeG,
	"H": KeyCodeH,
	"I": KeyCodeI,
	"J": KeyCodeJ,
	"K": KeyCodeK,
	"L": KeyCodeL,
	"M": KeyCodeM,
	"N": KeyCodeN,
	"O": KeyCodeO,
	"P": KeyCodeP,
	"Q": KeyCodeQ,
	"R": KeyCodeR,
	"S": KeyCodeS,
	"T": KeyCodeT,
	"U": KeyCodeU,
	"V": KeyCodeV,
	"W": KeyCodeW,
	"X": KeyCodeX,
	"Y": KeyCodeY,
	"Z": KeyCodeZ,

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

	"Return": KeyCodeReturnEnter,
	"Tab":    KeyCodeTab,
	"Space":  KeyCodeSpacebar,
	"-":      KeyCodeHyphenMinus,        // -
	"=":      KeyCodeEqualSign,          // =
	"[":      KeyCodeLeftSquareBracket,  // [
	"]":      KeyCodeRightSquareBracket, // ]
	"\\":     KeyCodeBackslash,          // \
	";":      KeyCodeSemicolon,          // ;
	"'":      KeyCodeApostrophe,         // '
	"`":      KeyCodeGraveAccent,        // `
	",":      KeyCodeComma,              // ,
	".":      KeyCodeFullStop,           // .
	"/":      KeyCodeSlash,              // /

	"Keypad /":     KeyCodeKeypadSlash,       // /
	"Keypad *":     KeyCodeKeypadAsterisk,    // *
	"Keypad -":     KeyCodeKeypadHyphenMinus, // -
	"Keypad +":     KeyCodeKeypadPlusSign,    // +
	"Keypad Enter": KeyCodeKeypadEnter,
	"Keypad 1":     KeyCodeKeypad1,
	"Keypad 2":     KeyCodeKeypad2,
	"Keypad 3":     KeyCodeKeypad3,
	"Keypad 4":     KeyCodeKeypad4,
	"Keypad 5":     KeyCodeKeypad5,
	"Keypad 6":     KeyCodeKeypad6,
	"Keypad 7":     KeyCodeKeypad7,
	"Keypad 8":     KeyCodeKeypad8,
	"Keypad 9":     KeyCodeKeypad9,
	"Keypad 0":     KeyCodeKeypad0,
	"Keypad .":     KeyCodeKeypadFullStop,  // .
	"Keypad =":     KeyCodeKeypadEqualSign, // =

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

	"Right": KeyCodeRightArrow,
	"Left":  KeyCodeLeftArrow,
	"Down":  KeyCodeDownArrow,
	"Up":    KeyCodeUpArrow,

	"Numlock": KeyCodeKeypadNumLock,

	"Help": KeyCodeHelp,

	"AudioMute":  KeyCodeMute,
	"Mute":       KeyCodeMute,
	"VolumeUp":   KeyCodeVolumeUp,
	"VolumeDown": KeyCodeVolumeDown,

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

	"Left Ctrl":     KeyCodeLeftControl,
	"Left Shift":    KeyCodeLeftShift,
	"Left Alt":      KeyCodeLeftAlt,
	"Left Option":   KeyCodeLeftAlt,
	"Left GUI":      KeyCodeLeftGUI,
	"Left Command":  KeyCodeLeftGUI,
	"Right Ctrl":    KeyCodeRightControl,
	"Right Shift":   KeyCodeRightShift,
	"Right Alt":     KeyCodeRightAlt,
	"Right Option":  KeyCodeLeftAlt,
	"Right GUI":     KeyCodeRightGUI,
	"Right Command": KeyCodeLeftGUI,

	// Compose

	//"": KeyCodeCompose,
}
