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
	OnPause()
	OnStop()
	OnDispose() error
}

// Runtime API
type Runtime interface {
	Use(plugin Plugin)
	Stop()
}

// Plugin API
type Plugin interface {
	Init(runtime Runtime) error
	Dispose()
}
