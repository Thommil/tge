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
	Key  string
	Type Type
}
