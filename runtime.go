package tge

import (
	sync "sync"
	time "time"
)

// App defines API to implement for TGE applications
type App interface {
	OnCreate(settings *Settings) error
	OnStart(runtime Runtime) error
	OnResize(width int, height int)
	OnResume()
	OnRender(elaspedTime time.Duration, mutex *sync.Mutex)
	OnTick(elaspedTime time.Duration, mutex *sync.Mutex)
	OnMouseEvent(event MouseEvent)
	OnScrollEvent(event ScrollEvent)
	OnKeyEvent(event KeyEvent)
	OnPause()
	OnStop()
	OnDispose() error
}

// Runtime API
type Runtime interface {
	Use(plugin Plugin)
	GetPlugin(name string) Plugin
	GetHost() interface{}
	GetRenderer() interface{}
	Stop()
}

// Plugin API
type Plugin interface {
	Init(runtime Runtime) error
	GetName() string
	Dispose()
}

// Type used in events to indicate type of action
type Type byte

// Button used in MouseEvent to indicate button
type Button byte

// KeyCode based on gomobile ones, sued for key mapping
type KeyCode int

const (
	// ButtonNone Button for not available or not applicable
	ButtonNone Button = 0
	// ButtonLeft Button for left button
	ButtonLeft Button = 1
	// ButtonMiddle Button for middle button
	ButtonMiddle Button = 2
	// ButtonRight Button for right button
	ButtonRight Button = 3
)

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

// MouseEvent definition
type MouseEvent struct {
	X, Y   int32
	Button Button
	Type   Type
}

// ScrollEvent definition
type ScrollEvent struct {
	X, Y int32
}

// KeyEvent definition
type KeyEvent struct {
	Key   KeyCode
	Value string
	Type  Type
}

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

// IsValid indicates a supported KeyCode
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
