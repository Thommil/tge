# TGE - Portable Runtime in GO
**TGE** aims to provide a light, portable and unopiniated runtime to integrate your favorite Go libraries (at least mines) to portable applications (dektop, web & mobile).

**TGE** is not and should not be another new game engine but instead a way to focus on business code and not low level pipes and hacks. The core is intended to be as light as possible and depends on plugins to enable cool features (OpenGL, AL, Vulkan, GUI...)

See it in action in [tge-examples](https://github.com/thommil/tge-examples) and look below for more details.

An online demo (Work in Progress) is also available here : [http://tge-demo.thommil.com](http://tge-demo.thommil.com)

## TGE Core 
Core implements a minimal runtime for several platforms:
  * Windows, MacOS & Linux
  * Android 5+ & IOS 8+
  * Web on Chrome, Firefox & Safari


<div style="text-align:center">
<img src="https://raw.githubusercontent.com/thommil/tge/master/specs/api.png"/>
</div>

## TGE plugins
Plugins allow to create your game/application by choosing implementation for different parts (GUI, rendering, sound...). This way, **you choose** how to create your game/application. I'm trying to port as many features as I can in plugins, natively portable libraries are off course directly usable too (physics, AI...).

Plugin&nbsp;link | Details
------------ | -------------
[tge-gl](https://github.com/thommil/tge-gl) | OpenGL/ES 3+ API based on [go-gl](https://github.com/go-gl/gl) work. Expose OpenGL to App in a portable way, il you want a low level access to the graphical context. Needed by most graphical libraries and plugins.
[tge-g3n](https://github.com/thommil/tge-g3n) | Based on the awesome [G3N](https://github.com/g3n/engine) game engine written in Go. This great piece of software engineering is brought to Mobile and Web by TGE.

# Getting started
Based on what is done with tools like Vue.js, Node or Spring, TGE offers a command line tool [tge-cli](https://github.com/thommil/tge-cli) to ease creation and build of TGE applications.

## Install tge-cli
To get the client, run:
```shell
go get github.com/thommil/tge-cli
```

The command line tool should be available in the $GOPATH/bin folder.

## Create new application
The create a new application workspace, run:
```shell
tge-cli init [package-name]
```

An application folder will be created with all needed resources to begin. See [tge-cli](https://github.com/thommil/tge-cli) and [Go Doc](https://godoc.org/github.com/thommil/tge) for details and API.

## Build the application
Once the application folder is created, releases can be generated using:
```shell
tge-cli build -target [target] [package-path]
```
Target allows to build your application for Desktop, Mobile or Web backend. See [tge-cli](https://github.com/thommil/tge-cli) for full details on how to use it and customize each target.

# Coding
## Applications
WIP 

The main entry point

```golang
package main

import (
	tge "github.com/thommil/tge"
  
)

type MyApp struct {
}

func (app *MyApp) OnCreate(settings *tge.Settings) error {
	return nil
}

func (app *MyApp) OnStart(runtime tge.Runtime) error {
	return nil
}

func (app *MyApp) OnResume() {
}

func (app *MyApp) OnRender(elapsedTime time.Duration, mutex *sync.Mutex) {
}

func (app *MyApp) OnTick(elapsedTime time.Duration, mutex *sync.Mutex) {
}

func (app *MyApp) OnPause() {
}

func (app *MyApp) OnStop() {
}

func (app *MyApp) OnDispose() error {
	return nil
}

func main() {
	tge.Run(&MyApp{})
}
```

## Plugins
WIP 

It's also possible to create new TGE plugins for sharing or to defines dedicated libraris accross your applications.

The paradigm of plugins is really simple, the code below:

```golang
package myplugin

import (
	tge "github.com/thommil/tge"
)

type plugin struct {
}

func init() {
	tge.Register(plugin{})
}

func (p *plugin) Init(runtime tge.Runtime) error {
	return nil
}

func (p *plugin) GetName() string {
	return "myplugin"
}

func (p *plugin) Dispose() {
}
```

## Targeting platform
WIP 

To define dedicated code to a specific platform, use the ***build*** preprocessing instruction of Go: