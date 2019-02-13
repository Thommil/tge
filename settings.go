package tge

import (
	physics "github.com/thommil/tge/physics"
	player "github.com/thommil/tge/player"
	renderer "github.com/thommil/tge/renderer"
	ui "github.com/thommil/tge/ui"
)

// Settings definition of TGE application
type Settings struct {
	// Name of the App
	Name string `json:"name" yaml:"name"`
	// Fullscreen indicates if the app must be run in fullscreen mode
	Fullscreen bool `json:"fullscreen" yaml:"fullscreen"`
	// Width of the window if run windowed only
	Width int `json:"width" yaml:"width"`
	// Height of the window if run windowed only
	Height int `json:"height" yaml:"height"`
	// Physics settings
	Physics physics.Settings `json:"physics" yaml:"physics"`
	// Player settings
	Player player.Settings `json:"player" yaml:"player"`
	// Renderer settings
	Renderer renderer.Settings `json:"renderer" yaml:"renderer"`
	// UI settings
	UI ui.Settings `json:"ui" yaml:"ui"`
}

// Default settings
var defaultSettings = Settings{
	Name:       "TGE Application",
	Fullscreen: false,
	Width:      640,
	Height:     480,
	Physics: physics.Settings{
		TPS: 100,
	},
	Player: player.Settings{},
	Renderer: renderer.Settings{
		FPS: 60,
	},
	UI: ui.Settings{},
}
