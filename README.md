# TGE - Portable Runtime in GO
**TGE** aims to provide a light, portable and unopiniated runtime to integrate your favorite GO libraries (at least mines) to portable applications (dektop, web & mobile).

**TGE** is not and should not be another new game engine but instead a way to focus on business code and not low level pipes and hacks. 

See it in action in [tge-examples](https://github.com/thommil/tge-examples) and look below for more details.

## TGE Core 
Core implements a minimal runtime for several platforms:
  * Windows, MacOS & Linux
  * Android 5+ & IOS 8+
  * Web on Chrome, Firefox & Safari

For diagrams lovers, here's the main architecture:

<p style="text-align:center">
<img src="https://raw.githubusercontent.com/thommil/tge/master/specs/api.png"/>
</p>

If you're wondering why App has been defined like this or what it the differences between Render() and Tick(), go directly to [Getting started](#getting-started) section and associated [Go Doc](https://godoc.org/github.com/thommil/tge).

## TGE plugins
Plugins then allow to create your game/application by choosing implementation of different parts (GUI, rendering, sound...). If a library needs a custom implementation to run on each target, I try to create the associated plugin, if not, just import it by yourself (ex: physics, AI...). This way, **you choose** how to create your game/application. 

Plugin&nbsp;link | Details
------------ | -------------
[tge-gl](https://github.com/thommil/tge-gl) | OpenGL/ES 3+ API based on [go-gl](https://github.com/go-gl/gl) work. Expose OpenGL to App in a portable way, il you want a low level access to the graphical context and needed by most graphical libraries and plugins.
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

Ideally, the package-name should be based on standard Go package rules (ex: github.com/thommil/my-app) but local package also works (ex: my-app).

An application folder will be created with all needed resources to begin. See [Go Doc](https://godoc.org/github.com/thommil/tge) for details.

## Build the application
Once the application folder is created, releases can be generated using:
```shell
tge-cli build -target [target] [package-path]
```
Target allows to build your application for Desktop, Mobile or Web backend. See [tge-cli](https://github.com/thommil/tge-cli) for full details on how to use it.
