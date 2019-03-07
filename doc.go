// Copyright (c) 2019 Thomas MILLET. All rights reserved.

/*
TGE Core contains interfaces and core implementation for supported targets:
 - desktop : MacOS, Linux, Windows
 - android : Android 5+
 - ios     : IOS 8+ (Work in progress)
 - browser : Chrome, Firefox, Safari (limited support)

TGE Core should not be used directly, it only defines interfaces and is used
by TGE Command Line Tool :
 see https://github.com/thommil/tge-cli

App

An App is the main entry point of TGE, the main() function should normally just
starts the Runtime, any other code not handled by the Runtime is potenitally
not portable:

 import "github.com/thommil/tge"

 func main() {
	tge.Run(&MyApp{})
 }

The App interface is described here and the implementation details in the auto
generated app.go using tge-cli.

Runtime

The Runtime instance is initialized through the Run(*App) function of main package. At
startup, the Runtime looks for registered plugins and initializes them. Then the
App instance is initialized and started.

The Runtime instance also exposes API for loading assets and subscribing to events
in a generic way.

Runtime exposes none portable objects like Host (backend) and Renderer (graphical context),
they can be used to implement custom behaviour depending on target in Apps or Plugins, the
implementations are as follows:

 Host:
  - desktop      : *sdl.Window - SDL2 from https://github.com/veandco/go-sdl2
  - android/ios  : mobile.App  - Custom gomobile from https://github.com/thommil/tge-mobile
  - browser      : *js.Value   - Gobal element through WebAssembly from Go 1.12

 Renderer:
  - desktop      : *sdl.GLContext - SDL2 from https://github.com/veandco/go-sdl2
  - android/ios  : gl.Context     - Custom gomobile from https://github.com/thommil/tge-mobile
  - browser      : *js.Value      - WebGL/WebGL2 context through WebAssembly from Go 1.12

Rendering

TGE uses Go channel mechanism to handle rendering, two loops are running side by side:

 - Ticker loop with Tick()   : handle CPU treatments (physics, AI, logical) and trigger rendering
 - Render loop with Render() : handle GPU treatments (draw calls)

Both loops are synchronized using a dedicated channel passed in parameter of each method. As this method
allows to make CPU/GPU treatments asynchronous, shared objects between contexts must correctly handled to avoid
conflicts. The sync channel is typed as interface{}, it can also be used to pass content and select specific treatments
based on underlying interface type. See examples for more details.

Events

Minimal set of events is handled by Runtime at the most possible portable way. Events
are then propagated through publish/subscribe:

 Subscribe(channel string, listener Listener)
 Unsubscribe(channel string, listener Listener)
 Publish(event Event)

Events are in their raw form (ie modifiers or gestures are not handled). It's up to the
application to implement specific needs. The aim of this approach is to keep the runtime
generic and fast by limiting treatments.

A dedicated plugin to generate advanced events will be available soon.

Plugins

As TGE core is intended to be as light as possible, all heavy treatments are deported to
plugins. The goal is to offer a portable API from Plugins by relying on Runtime.

Plugins are automatically registered at Go init() step, to use it, just import them as
standard Go packages, ex:

 import "github.com/thommil/tge-gl"

 func (app *App) OnStart(runtime tge.Runtime) error{
	 gl.ClearColor(0, 0, 0, 1)
 }

It's also possible to create custom plugins by implementing Plugin interface and
registering it in the Go init() function :

 package myplugin

 import (
	tge "github.com/thommil/tge"
 )

 type plugin struct {
	 Name string
 }

 func init() {
 	tge.Register(plugin{"myplugin"})
 }

 func (p *plugin) Init(runtime tge.Runtime) error {
	// Init code HERE if needed
	return nil
 }

 func (p *plugin) GetName() string {
 	return p.Name
 }

 func (p *plugin) Dispose() {
	 // Dispose code HERE if needed
 }

Targeting platform and Debug Mode

It's possible to write code for a specific platform the same way TGE do it.

For desktop, add the following Go build directives:

 // +build darwin freebsd linux windows
 // +build !android
 // +build !ios
 // +build !js

for mobile:

 // +build android ios

and for browser:

 // +build js

At last, it's also possible to create a dedicated file for debugging purpose by
adding:

 // +build debug

The file will be used if the -debug flag is set in tge-cli command line for build.

*/
package tge // import "github.com/thommil/tge"
