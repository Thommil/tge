package tge

// EventMask defines mask event for enable/disable events receivers
type EventMask int

const (
	// AllEventsDisable disables all input events on App
	AllEventsDisable = 0x00
	// MouseEventEnabled enabled mouse event receiver on App
	MouseEventEnabled = 0x01
	// ScrollEventEnabled enabled scroll event receiver on App
	ScrollEventEnabled = 0x02
	// KeyEventEnabled enabled key event receiver on App
	KeyEventEnabled = 0x04
	// AllEventsEnabled enables all input events on App
	AllEventsEnabled = MouseEventEnabled | ScrollEventEnabled | KeyEventEnabled
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
	// FPS (frames per seconds) target rate
	FPS int `json:"fps" yaml:"fps"`
	// TPS (ticks per seconds) target rate
	TPS int `json:"tps" yaml:"tps"`
	// EventMask allows to enabled/disable events receiver on App
	EventMask EventMask `json:"event_mask" yaml:"event_mask"`
}

// Default settings
var defaultSettings = Settings{
	Name:       "TGE Application",
	Fullscreen: false,
	Width:      640,
	Height:     480,
	TPS:        100,
	FPS:        60,
	EventMask:  AllEventsEnabled,
}
