package tge

import (
	context "github.com/thommil/tge/context"
	physics "github.com/thommil/tge/physics"
	sound "github.com/thommil/tge/sound"
	ui "github.com/thommil/tge/ui"
)

// Runtime API
type Runtime interface {
	GetContext() context.Context
	GetUI() ui.UI
	GetPhysics() physics.Physics
	GetSound() sound.Sound
}
