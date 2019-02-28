# TGE - Game Runtime in GO
**TGE** aims to provide a light, portable and almsot unopiniated runtime to integrate your favorite GO libraries (at least mines).

## TGE Core 
Core implements a minimal runtime for several platforms:
  * Windows, MacOS & Linux based on [SDL2](https://github.com/veandco/go-sdl2)
  * Android 5+ & IOS 8+ based on [tge-mobile](https://github.com/thommil/tge-mobile)
  * Web on Chrome, Firefox & Safari based on [Go WebAssemby](https://github.com/golang/go/wiki/WebAssembly)

For diagrams lovers, here's the main architecture:

![Core Architecture](https://raw.githubusercontent.com/thommil/tge/master/specs/api.png)

## TGE plugins
Plugins then allow to create your game/application by choosing implementation of different parts (GUI, rendering, sound, physics ...). If a library needs a custom implementation to run on each target, I try to create the associated plugin, if not, just import it by yourself. This way, **you choose** how to create your game/application. 

Currently available plugins :
 * [tge-gl](https://github.com/thommil/tge-gl) : OpenGL/ES 3+ API 
 * [tge-g3n](https://github.com/thommil/tge-g3n) : Port of the awesome [G3N](https://github.com/g3n/engine) Game Engine


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

Ideally, the package-name should be based on standard Go package rules (ex: gihtub.com/thommil/my-app) but local package also works (ex: my-app).


