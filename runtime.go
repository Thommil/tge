// Copyright (c) 2019 Thomas MILLET. All rights reserved.

package tge

import (
	fmt "fmt"
	reflect "reflect"
	sync "sync"
	time "time"
)

// Runtime singleton
var _runtimeInstance Runtime

// -------------------------------------------------------------------- //
// API
// -------------------------------------------------------------------- //

// App is the main entry point of a TGE Application. Created using TGE Command Line Tools.
//
// See generated app.go by tge-cli for more details and explanations.
type App interface {
	// OnCreate is called at App instanciation, the Runtime resources are not
	// yet available. Settings and treatments not related to Runtime should be done here.
	OnCreate(settings *Settings) error

	// OnStart is called when all Runtime resources are available but looping as not been
	// started. Initializations should be done here (GL conf, physics engine ...).
	OnStart(runtime Runtime) error

	// OnResume is called just after OnStart and also when the Runtime is awaken after a pause.
	OnResume()

	// OnRender is called each time the graphical context is redrawn. This method should only implements
	// graphical calls and not logical ones. The mutex allows to synchronize critical path with Tick() loop.
	OnRender(elaspedTime time.Duration, mutex *sync.Mutex)

	// OnTick is called at a rate defined in settings, logical operations should be done here like physics, AI
	// or any background task not relatd to graphics. The mutex allows to synchronize critical path with Render() loop.
	OnTick(elaspedTime time.Duration, mutex *sync.Mutex)

	// OnPause is called when the Runtime lose focus (alt-tab, home button, tab change ...) This is a good nrty point to
	// set and display a pause screen
	OnPause()

	// OnStop is called when the Runtime is ending, context saving should be done here. On current Android version this
	// handler is also called when the application is paused (and restart after).
	OnStop()

	// OnDispose is called when all exit treatments are done for cleaning task (memory, tmp files ...)
	OnDispose()
}

// Listener is the callback definition for publish/subscribe, the return value indicates if the event has been consumed (true)
// and propagation stopped, in other case the next registered Listener is called
type Listener func(event Event) bool

// Runtime defines the commmon API across runtimes implementations
type Runtime interface {
	// GetAsset retrieves assets in []byte form, assets are always stored in
	// package asset folder independently of the target
	GetAsset(path string) ([]byte, error)

	// GetHost is for low level and target specific implementation and allows to
	// retrieve the underlying backend of the Runtime (see package description)
	GetHost() interface{}

	// GetRenderer is for low level and target specific implementation and allows to
	// retrieve the underlying graphical context of the Runtime (see package description)
	GetRenderer() interface{}

	// Subscribe register a new Listener to specified channel
	Subscribe(channel string, listener Listener)

	// Unsubscribe deregister a new Listener from specified channel
	Unsubscribe(channel string, listener Listener)

	// Publish send an Event on channel defined in the Event.Channel()
	Publish(event Event)

	// Stop allows App to end the Runtime directly
	Stop()
}

// -------------------------------------------------------------------- //
// Plugins
// -------------------------------------------------------------------- //

// Plugin interface defines the API used by the Runtime to handle plugins
type Plugin interface {
	// Init is called before Runtime looping
	Init(runtime Runtime) error

	// GetName allows Runtime to identify plugin by its name
	GetName() string

	// Dispose is called at end of Runtime before exiting
	Dispose()
}

// Inner map of plugins
var plugins = make(map[string]Plugin)

// Register a plugin in Runtime, this function should only be called
// in the Go init() function of plugins to allow registration of plugins
// before looping. In other case, the Init() method of plugins are never called.
func Register(plugin Plugin) {
	name := plugin.GetName()
	if _, found := plugins[name]; !found {
		plugins[name] = plugin
		fmt.Printf("Plugin %s registered\n", name)
	}
}

func initPlugins() {
	for _, plugin := range plugins {
		err := plugin.Init(_runtimeInstance)
		if err != nil {
			fmt.Printf("Failed to initialize plugin %s: %v\n", plugin.GetName(), err)
			panic(err)
		}
		fmt.Printf("Plugin %s loaded\n", plugin.GetName())
	}
}

// Global dispose
func dispose() {
	for _, plugin := range plugins {
		plugin.Dispose()
		fmt.Printf("Plugin %s released\n", plugin.GetName())
	}
}

// -------------------------------------------------------------------- //
// Events
// -------------------------------------------------------------------- //

// Type is a component of events to indicate a generic way of defining an
// event action.
type Type byte

// Button indicates the type of button or touch used in event
type Button byte

// KeyCode is used to map raw key codes values
type KeyCode int

// Buttons values
const (
	// ButtonNone Button for not available or not applicable
	ButtonNone Button = 0
	// ButtonLeft Button for left button, first finger touch
	ButtonLeft Button = 1
	// ButtonMiddle Button for middle button, second finger touch
	ButtonMiddle Button = 2
	// ButtonRight Button for right button, third finger touch
	ButtonRight Button = 3
)

// Types values
const (
	// TypeNone Type for not available or not applicable
	TypeNone Type = 0
	// TypeDown Type for pressed button/key/touch
	TypeDown Type = 1
	// TypeUp Type for released button/key/touch
	TypeUp Type = 2
	// TypeMove Type for mouse/touch move
	TypeMove Type = 3
)

// Events interface defines an event base by its channel
type Event interface {
	// Type defines a unique keywork/channel for event
	Channel() string
}

// ResizeEvent is triggered when TGE painting area is resized
type ResizeEvent struct {
	Width, Height int32
}

// Channel of ResizeEvent = "resize"
func (e ResizeEvent) Channel() string {
	return "resize"
}

// MouseEvent is triggered on mouse/touch down/up event and
// mouse motion event too
type MouseEvent struct {
	X, Y   int32
	Button Button
	Type   Type
}

// Channel of MouseEvent = "mouse"
func (e MouseEvent) Channel() string {
	return "mouse"
}

// ScrollEvent is called only on desktop/browser, X/Y values are
// only [-1, 0, 1] to normalize scrolling accross targets
type ScrollEvent struct {
	X, Y int32
}

// Channel of ScrollEvent = "scroll"
func (e ScrollEvent) Channel() string {
	return "scroll"
}

// KeyEvent defines a down/up key event, the Key attribute is portable
// accross targets, the Value is the string representation of the key
type KeyEvent struct {
	Key   KeyCode
	Value string
	Type  Type
}

// Channel of KeyEvent = "key"
func (e KeyEvent) Channel() string {
	return "key"
}

// Inner map of listeners
var listeners = make(map[string][]Listener)

func subscribe(channel string, listener Listener) {
	if _, found := listeners[channel]; !found {
		listeners[channel] = make([]Listener, 0, 10)
	}
	listeners[channel] = append(listeners[channel], listener)
}

func unsubscribe(channel string, listener Listener) {
	if all, found := listeners[channel]; found {
		for i, l := range all {
			if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
				listeners[channel] = append(listeners[channel][:i], listeners[channel][i+1:]...)
				break
			}
		}
	}
}

func publish(event Event) {
	if list, found := listeners[event.Channel()]; found {
		for _, listener := range list {
			if listener(event) {
				break
			}
		}
	}
}

// Keycode constants
const (
	// Unkwown

	KeyCodeUnknown KeyCode = 0

	// Printable
	KeyCodeA KeyCode = 1
	KeyCodeB KeyCode = 2
	KeyCodeC KeyCode = 3
	KeyCodeD KeyCode = 4
	KeyCodeE KeyCode = 5
	KeyCodeF KeyCode = 6
	KeyCodeG KeyCode = 7
	KeyCodeH KeyCode = 8
	KeyCodeI KeyCode = 9
	KeyCodeJ KeyCode = 10
	KeyCodeK KeyCode = 11
	KeyCodeL KeyCode = 12
	KeyCodeM KeyCode = 13
	KeyCodeN KeyCode = 14
	KeyCodeO KeyCode = 15
	KeyCodeP KeyCode = 16
	KeyCodeQ KeyCode = 17
	KeyCodeR KeyCode = 18
	KeyCodeS KeyCode = 19
	KeyCodeT KeyCode = 20
	KeyCodeU KeyCode = 21
	KeyCodeV KeyCode = 22
	KeyCodeW KeyCode = 23
	KeyCodeX KeyCode = 24
	KeyCodeY KeyCode = 25
	KeyCodeZ KeyCode = 26

	KeyCode1 KeyCode = 27
	KeyCode2 KeyCode = 28
	KeyCode3 KeyCode = 29
	KeyCode4 KeyCode = 30
	KeyCode5 KeyCode = 31
	KeyCode6 KeyCode = 32
	KeyCode7 KeyCode = 33
	KeyCode8 KeyCode = 34
	KeyCode9 KeyCode = 35
	KeyCode0 KeyCode = 36

	KeyCodeReturnEnter        KeyCode = 37
	KeyCodeTab                KeyCode = 38
	KeyCodeSpacebar           KeyCode = 39
	KeyCodeHyphenMinus        KeyCode = 40 // -
	KeyCodeEqualSign          KeyCode = 41 // =
	KeyCodeLeftSquareBracket  KeyCode = 42 // [
	KeyCodeRightSquareBracket KeyCode = 43 // ]
	KeyCodeBackslash          KeyCode = 44 // \
	KeyCodeSemicolon          KeyCode = 45 // ;
	KeyCodeApostrophe         KeyCode = 46 // '
	KeyCodeGraveAccent        KeyCode = 47 // `
	KeyCodeComma              KeyCode = 48 // ,
	KeyCodeFullStop           KeyCode = 49 // .
	KeyCodeSlash              KeyCode = 50 // /

	KeyCodeKeypadSlash       KeyCode = 51 // /
	KeyCodeKeypadAsterisk    KeyCode = 52 // *
	KeyCodeKeypadHyphenMinus KeyCode = 53 // -
	KeyCodeKeypadPlusSign    KeyCode = 54 // +
	KeyCodeKeypadEnter       KeyCode = 55
	KeyCodeKeypad1           KeyCode = 56
	KeyCodeKeypad2           KeyCode = 57
	KeyCodeKeypad3           KeyCode = 58
	KeyCodeKeypad4           KeyCode = 59
	KeyCodeKeypad5           KeyCode = 60
	KeyCodeKeypad6           KeyCode = 61
	KeyCodeKeypad7           KeyCode = 62
	KeyCodeKeypad8           KeyCode = 63
	KeyCodeKeypad9           KeyCode = 64
	KeyCodeKeypad0           KeyCode = 65
	KeyCodeKeypadFullStop    KeyCode = 66 // .
	KeyCodeKeypadEqualSign   KeyCode = 67 // =

	KeyCodeAt                KeyCode = 68 // @
	KeyCodeGreaterThan       KeyCode = 69 // >
	KeyCodeLesserThan        KeyCode = 70 // <
	KeyCodeDollar            KeyCode = 71 // $
	KeyCodeColon             KeyCode = 72 // :
	KeyCodeLeftParenthesis   KeyCode = 73 // (
	KeyCodeLRightParenthesis KeyCode = 74 // )

	KeyCodeAmpersand   KeyCode = 75 // &
	KeyCodeHash        KeyCode = 76 // #
	KeyDoubleQuote     KeyCode = 77 // "
	KeyQuote           KeyCode = 78 // '
	KeyParapgrah       KeyCode = 79 // §
	KeyExclamationMark KeyCode = 80 // !
	KeyUnderscore      KeyCode = 81 // _
	KeyQuestionMark    KeyCode = 82 // ?
	KeyPercent         KeyCode = 83 // %
	KeyDegree          KeyCode = 84 // °

	// Actions

	KeyCodeEscape   KeyCode = 101
	KeyCodeCapsLock KeyCode = 102

	KeyCodeDeleteBackspace KeyCode = 103
	KeyCodePause           KeyCode = 104
	KeyCodeInsert          KeyCode = 105
	KeyCodeHome            KeyCode = 106
	KeyCodePageUp          KeyCode = 107
	KeyCodeDeleteForward   KeyCode = 108
	KeyCodeEnd             KeyCode = 109
	KeyCodePageDown        KeyCode = 110

	KeyCodeRightArrow KeyCode = 111
	KeyCodeLeftArrow  KeyCode = 112
	KeyCodeDownArrow  KeyCode = 113
	KeyCodeUpArrow    KeyCode = 114

	KeyCodeKeypadNumLock KeyCode = 115

	KeyCodeHelp KeyCode = 116

	KeyCodeMute       KeyCode = 120
	KeyCodeVolumeUp   KeyCode = 121
	KeyCodeVolumeDown KeyCode = 122

	// Functions

	KeyCodeF1  KeyCode = 201
	KeyCodeF2  KeyCode = 202
	KeyCodeF3  KeyCode = 203
	KeyCodeF4  KeyCode = 204
	KeyCodeF5  KeyCode = 205
	KeyCodeF6  KeyCode = 206
	KeyCodeF7  KeyCode = 207
	KeyCodeF8  KeyCode = 208
	KeyCodeF9  KeyCode = 209
	KeyCodeF10 KeyCode = 210
	KeyCodeF11 KeyCode = 211
	KeyCodeF12 KeyCode = 212

	// Modifiers

	KeyCodeLeftControl  KeyCode = 301
	KeyCodeLeftShift    KeyCode = 302
	KeyCodeLeftAlt      KeyCode = 303
	KeyCodeLeftGUI      KeyCode = 304
	KeyCodeRightControl KeyCode = 305
	KeyCodeRightShift   KeyCode = 306
	KeyCodeRightAlt     KeyCode = 307
	KeyCodeRightGUI     KeyCode = 308

	// Compose

	KeyCodeCompose KeyCode = 0x10000
)

// IsValid indicates a recongonized/valid KeyCode
func (k KeyCode) IsValid() bool {
	return k != KeyCodeUnknown
}

// IsPrintable indicates a printable KeyCode
func (k KeyCode) IsPrintable() bool {
	return (k >= KeyCodeUnknown) && (k < KeyCodeEscape)
}

// IsAction indicates an action KeyCode
func (k KeyCode) IsAction() bool {
	return (k >= KeyCodeEscape) && (k < KeyCodeF1)
}

// IsFunction indicates a function KeyCode
func (k KeyCode) IsFunction() bool {
	return (k >= KeyCodeF1) && (k < KeyCodeLeftControl)
}

// IsModifier indicates a modifier KeyCode
func (k KeyCode) IsModifier() bool {
	return (k >= KeyCodeLeftControl) && (k < KeyCodeCompose)
}

// IsCompose indicates a comopose KeyCode
func (k KeyCode) IsCompose() bool {
	return k == KeyCodeCompose
}
