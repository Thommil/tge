package tge

import (
	physics "github.com/thommil/tge/physics"
	player "github.com/thommil/tge/player"
	renderer "github.com/thommil/tge/renderer"
	ui "github.com/thommil/tge/ui"
)

// App defines API to implement for TGE applications
type App interface {
	Create(settings Settings) error
	Start(runtime Runtime) error
	Resize(width int, height int)
	Resume()
	Render(renderer renderer.Renderer, ui ui.UI, player player.Player)
	Tick(physics physics.Physics)
	Pause()
	Dispose() error
}

// Runtime API
type Runtime interface {
}

// Instanciate is the main entry point
func Instanciate(app App) {
	app.Create(Settings{})

	err := backend_Instanciate(app)

	if err != nil {

	}
	//app.Start()
	//app.Resume()
	//app.Resize()

	//app.Render()
	//app.Tick()

	//app.Dispose()
}

type BasicApp struct {
}

func (app BasicApp) Create(settings tge.Settings) error {
	log.Println("Create()")
	return nil
}

func (app BasicApp) Start(runtime tge.Runtime) error {
	log.Println("Start()")
	return nil
}

func (app BasicApp) Resize(width int, height int) {
	log.Println("Resize()")
}

func (app BasicApp) Resume() {
	log.Println("Resume()")
}

func (app BasicApp) Render(renderer renderer.Renderer, ui ui.UI, player player.Player) {
	log.Println("Render()")
}

func (app BasicApp) Tick(physics physics.Physics) {
	log.Println("Tick()")
}

func (app BasicApp) Pause() {
	log.Println("Pause()")
}

func (app BasicApp) Dispose() error {
	log.Println("Dispose()")
	return nil
}

func main() {
	app := BasicApp{}
	tge.Instanciate(app)
}
