package tge

import (
	log "log"
	sync "sync"
	time "time"
)

// App defines API to implement for TGE applications
type App interface {
	OnCreate(settings *Settings) error
	OnStart(runtime Runtime) error
	OnResize(width int, height int)
	OnResume()
	OnRender(elaspedTime time.Duration, locker sync.Locker)
	OnTick(elaspedTime time.Duration, locker sync.Locker)
	OnPause()
	OnStop()
	OnDispose() error
}

// Runtime API
type Runtime interface {
	Stop()
}

func init() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
}

// Run is the main entry point
func Run(app App) {
	log.Println("Run()")

	settings := &defaultSettings
	app.OnCreate(settings)

	err := doRun(app, settings)
	if err != nil {
		log.Fatalln(err)
	}
	defer app.OnDispose()
}
