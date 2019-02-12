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
}

// Default settings
var defaultSettings = Settings{
	Name:       "Untilted",
	Fullscreen: true,
	Width:      640,
	Height:     480,
}
