package tge

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
	// FPS (frames per seconds) target rate
	FPS int `json:"fps" yaml:"fps"`
	// TPS (ticks per seconds) target rate
	TPS int `json:"tps" yaml:"tps"`
}

// Default settings
var defaultSettings = Settings{
	Name:       "TGE Application",
	Fullscreen: false,
	Width:      640,
	Height:     480,
	TPS:        100,
	FPS:        60,
}
